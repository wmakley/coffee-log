package middleware

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthMiddlewareOptions struct {
	DbConn      *sql.DB
	Realm       string
	MaxAttempts int32
	Debug       bool
}

var ErrNotAUser = errors.New("user is not *sqlc.User")

// GetCurrentUser gets the current user from the gin context
// If the user does not exist, user will be nil and exists
// will be false. If user is not nil, but cannot be
// converted to *User, it will panic.
func GetCurrentUser(c *gin.Context) (*sqlc.User, bool) {
	value, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	user, ok := value.(*sqlc.User)
	if !ok {
		panic(ErrNotAUser)
	}
	return user, true
}

// AuthMiddleware returns a custom basic authentication middleware with fail2ban
// that uses the database to store bans and login attempts.
func AuthMiddleware(options AuthMiddlewareOptions) gin.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote(options.Realm)

	return func(c *gin.Context) {
		store := StoreFromCtx(c, options.DbConn)

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
			if options.Debug {
				log.Printf("authorization decode error: %+v", err)
			}
			authenticationFailure(c, realm)
			return
		}

		if options.Debug {
			log.Printf("got username: %s, password: %s from header", username, password)
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
	encodedCredentials := strings.TrimPrefix(rawHeader, "Basic ")

	usernamePasswdBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
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
