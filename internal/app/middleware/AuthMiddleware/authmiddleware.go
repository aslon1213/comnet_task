package authmiddleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	DB  *sql.DB
	ctx context.Context
}

func New(DB *sql.DB, ctx context.Context) *AuthMiddleware {
	return &AuthMiddleware{
		DB:  DB,
		ctx: ctx,
	}
}

// AuthMiddleware is a middleware which checks for valid token.
// it checks for SESSTOKEN cookie and it's value which is jwt token returned by /user/auth.
// if token is valid, it sets user_id and username to context and calls nextHandler by next().
func (a *AuthMiddleware) AuthMiddleware(c *gin.Context) {
	//get user token

	cookies := c.Request.Cookies()
	if len(cookies) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var token_cookie *http.Cookie
	// check whether SESSTOKEN cookie exists and value is not empty
	for _, cookie := range cookies {
		if cookie.Name == "SESSTOKEN" {
			token_cookie = cookie
			break
		}
	}
	if token_cookie == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token_cookie.Value == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenString := token_cookie.Value
	// check if token is valid
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret
		hmacSampleSecret := []byte(os.Getenv("SIGNING_SECRET"))
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check if token is expired or not
		expite_time_string, ok := claims["expires"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		expire_time, err := time.Parse(time.RFC3339, expite_time_string.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if expire_time.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			return
		}

		// if valid, get user id from token
		// set user_id and username to context and call next()
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
	}

}

// func(a *AuthMiddleware) CheckForRoutNames(c *gin.Context) {

// }
