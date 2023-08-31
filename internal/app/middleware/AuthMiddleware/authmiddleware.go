package authmiddleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

func (a *AuthMiddleware) AuthMiddleware(c *gin.Context) {
	//get user token

	cookies := c.Request.Cookies()
	// fmt.Println(cookies)
	if len(cookies) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var token_cookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			// fmt.Println(cookie.Value)
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
	// fmt.Println(tokenString)
	// check if token is valid
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SIGNING_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// if valid, get user id from token
		// fmt.Println(claims["username"], claims["expires"])
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
	}

	// if valid, get user from db

	// if user not found, return error

	// if user found, set user to context

	// next()
}
