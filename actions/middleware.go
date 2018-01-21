package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
	jose "gopkg.in/square/go-jose.v2"
)

// CurrentUserSetter sets current user info and roles in the context
func CurrentUserSetter(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		uid := c.Request().Header.Get("Uid")
		tx := c.Value("tx").(*pop.Connection)

		user := &models.User{}
		err := tx.Where("auth_id = ?", uid).First(user)

		if err != nil {
			return errors.WithStack(err)
		}

		roles := &models.Roles{}
		sql := "SELECT roles.* FROM user_roles INNER JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = ?"
		q := tx.RawQuery(sql, user.ID)
		err = q.All(roles)

		if err != nil {
			return errors.WithStack(err)
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

		err = next(c)
		return err
	}
}

// Authenticate will ensure only authenticated users gain access to protected endpoints
func Authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// do some work before calling the next handler
		fmt.Println(c.Request().Header["Authorization"])
		err := checkJwt(c.Response(), c.Request())
		if err == nil {
			err := next(c)
			// do some work after calling the next handler
			return err
		}

		return err
	}
}

func checkJwt(w http.ResponseWriter, r *http.Request) error {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: JwksURI})
	audience := Auth0APIAudience

	configuration := auth0.NewConfiguration(client, audience, Auth0APIIssuer, jose.RS256)
	validator := auth0.NewValidator(configuration)

	token, err := validator.ValidateRequest(r)

	if err != nil {
		fmt.Println(err.Error())

		response := Response{
			Message: "Missing or invalid token.",
		}

		fmt.Println(token)

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return err
	}

	return nil
}
