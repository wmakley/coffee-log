package internal

import (
	"coffee-log/db/sqlc"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddleOptions struct {
	Realm string
	MaxAttempts int32
	Debug bool
}

// AuthMiddleware returns a custom basic authentication middleware with fail2ban
// that uses the database to store bans and login attempts.
func AuthMiddleware(options AuthMiddleOptions) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote(options.Realm)

	return func(c *gin.Context) {
		store := sqlc.StoreFromCtx(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if options.Debug {
				log.Print("authorization header is empty")
			}
			authenticationFailure(c, realm)
			return
		}

		username, password, err := decodeUsernameAndPassword(authHeader)
		if err != nil {
			log.Printf("authorization decode error: %+v", err)
			authenticationFailure(c, realm)
			return
		}

		ip := c.ClientIP()
		user, err := store.CheckAndLogLoginAttempt(c, ip, username, password, options.MaxAttempts)
		if err != nil {
			if err == sqlc.ErrBadCredentials {
				if options.Debug {
					log.Print("bad username or password")
				}
				authenticationFailure(c, realm)
			} else if err == sqlc.ErrIPBanned {
				if options.Debug {
					log.Print("ip address is banned")
				}
				implementIPBan(c)
			} else {
				if options.Debug {
					log.Printf("unexpected error: %+v", err)
				}
				// roll back transaction
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		if options.Debug {
			log.Print("authentication success, username=", user.Username)
		}

		c.Set("user", &user)

		c.Next()
	}
}

func decodeUsernameAndPassword(rawHeader string) (username string, password string, err error) {
	encodedCredentials := []byte(strings.TrimPrefix(rawHeader, "Basic "))

	usernamePasswdBytes := make([]byte, len(encodedCredentials))
	if _, err = base64.StdEncoding.Decode(usernamePasswdBytes, encodedCredentials); err != nil {
		err = fmt.Errorf("base64 decode error: %w", err)
		return
	}

	usernamePasswdSplit := strings.SplitN(string(usernamePasswdBytes), ":", 2)
	if len(usernamePasswdSplit) != 2 {
		err = fmt.Errorf("invalid authorization header")
		return
	}

	username = usernamePasswdSplit[0]
	password = usernamePasswdSplit[1]
	return
}

func authenticationFailure(c *gin.Context, realm string) {
	c.Header("WWW-Authenticate", realm)
	c.AbortWithStatus(http.StatusUnauthorized)
}

func implementIPBan(c *gin.Context) {
	log.Printf("ip address %s is banned", c.ClientIP())
	c.AbortWithStatus(http.StatusNotFound)
}

func authorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}
