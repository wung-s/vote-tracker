package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
	"golang.org/x/net/context"
)

// Authenticate will ensure only authenticated users gain access to protected endpoints
func Authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler
		client, err := FirebaseApp.Auth(context.Background())

		idToken := c.Request().Header.Get("Authorization")
		idToken = strings.Replace(idToken, `bearer `, "", 1)
		if ENV == "development" {
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
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	err := tx.Where("auth_id = ?", uid).First(user)

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

	c.Set("currentUser", currentUser)
	return nil
}
