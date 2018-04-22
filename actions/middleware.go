package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
	"golang.org/x/net/context"
)

func oneWeek() time.Duration {
	return 7 * 24 * time.Hour
}

func RestrictedHandlerMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if len(tokenString) == 0 {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("No token set in headers"))
		}

		// parsing token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			box := packr.NewBox("../config")
			mySignedKey := box.Bytes(os.Getenv("JWT_SIGN_KEY"))

			return mySignedKey, nil
		})

		if err != nil {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Could not parse the token, %v", err))
		}

		// getting claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			setCurrentUser(claims["jti"].(string), c)
		} else {
			return c.Error(http.StatusUnauthorized, fmt.Errorf("Failed to validate token: %v", claims))
		}

		return next(c)
	}
}

// Authenticate will ensure only authenticated users gain access to protected endpoints
func Authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler
		client, err := FirebaseApp.Auth(context.Background())

		idToken := c.Request().Header.Get("Authorization")
		idToken = strings.Replace(idToken, `bearer `, "", 1)
		if ENV == "development" || ENV == "test" {
			fmt.Println("Authorization", idToken)
		}
		token, err := client.VerifyIDToken(idToken)
		if err != nil {
			fmt.Printf("error verifying ID token: %v\n", err)
			response := Response{
				Message: "Missing or invalid token.",
			}
			c.Response().WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(c.Response()).Encode(response)
			return err
		}

		if err := setCurrentUser(token.UID, c); err != nil {
			return errors.WithStack(err)
		}
		err = next(c)
		return err
	}
}

func setCurrentUser(uid string, c buffalo.Context) error {
	tx := models.DB
	user := &models.User{}
	err := tx.Where("id = ?", uid).First(user)

	if err != nil {
		return err
	}

	roles := &models.Roles{}
	sql := "SELECT roles.* FROM user_roles INNER JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = ?"
	q := tx.RawQuery(sql, user.ID)
	err = q.All(roles)

	if err != nil {
		return err
	}

	roleNames := []string{}
	// obtain only the name of the roles
	for _, v := range *roles {
		roleNames = append(roleNames, v.Name)
	}

	currentUser := struct {
		models.User
		Roles []string `json:"roles"`
	}{
		*user,
		roleNames,
	}

	currentUser.Password = ""
	c.Set("currentUser", currentUser)
	return nil
}
