package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/x/sessions"
	"github.com/wung-s/gotv/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// Auth0APIAudience identifies the server in Auth0
var Auth0APIAudience = []string{"https://gotv.com"}

const (
	// Auth0APIIssuer is the issuer
	Auth0APIIssuer = "https://wung.auth0.com/"
	JwksURI        = "https://wung.auth0.com/.well-known/jwks.json"
)

// Response defines the response message structure
type Response struct {
	Message string `json:"message"`
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_gotv_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		app.Use(Authenticate)
		app.Middleware.Skip(Authenticate, HomeHandler)

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		app.GET("/", HomeHandler)

		// app.Resource("/members", MembersResource{})
		app.POST("/members/upload", MembersUpload)
		// app.Resource("/roles", RolesResource{})
	}

	return app
}
