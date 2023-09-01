package utilshelpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aslon1213/comnet_task/internal/app/models"
	"github.com/golang-jwt/jwt/v5"
)

func Check_for_route_names(user models.User) error {
	name := strings.ToLower(user.Name)
	login := strings.ToLower(user.Login)

	if name == "register" || name == "auth" || name == "phone" {
		return fmt.Errorf("similar Name to route names not accepted")
	}
	if login == "register" || login == "auth" || login == "phone" {
		return fmt.Errorf("similar Login to route names not accepted")
	}
	return nil
}

func CreateSessionCookieToken(user models.User, expire_time time.Time) (string, error) {

	// one day for expiration
	// expirations_time := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Login,
		"user_id":  user.ID,
		"expires":  expire_time,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SIGNING_SECRET")))
	if err != nil {
		return "", err
	}

	// fmt.Println(tokenString)
	return tokenString, nil
}
