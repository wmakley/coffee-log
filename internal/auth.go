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

// AuthMiddleware returns a custom basic authentication middleware with fail2ban
// that uses the database to store bans and login attempts.
func AuthMiddleware(realm string, maxAttempts int32) gin.HandlerFunc {
	realm = "Basic realm=" + strconv.Quote(realm)

	return func(c *gin.Context) {
		store := sqlc.StoreFromCtx(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authenticationFailure(c, realm)
			return
		}

		username, password, err := decodeUsernameAndPassword(authHeader)
		if err != nil {
			log.Print("authorization decode error:", err.Error())
			authenticationFailure(c, realm)
			return
		}

		ip := c.ClientIP()
		user, err := store.CheckAndLogLoginAttempt(c, ip, username, password, maxAttempts)
		if err != nil {
			if err == sqlc.ErrBadCredentials {
				authenticationFailure(c, realm)
			} else if err == sqlc.ErrIPBanned {
				implementIPBan(c)
			} else {
				// roll back transaction
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.Set("user", &user)

		c.Next()
	}
}

func decodeUsernameAndPassword(rawHeader string) (username string, password string, err error) {
	encodedCredentials := strings.TrimPrefix(rawHeader, "Basic ")

	var usernamePasswdBytes []byte
	if _, err = base64.StdEncoding.Decode(usernamePasswdBytes, []byte(encodedCredentials)); err != nil {
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
