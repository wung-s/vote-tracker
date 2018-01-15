package actions

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"github.com/gobuffalo/buffalo"
	jose "gopkg.in/square/go-jose.v2"
)

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
