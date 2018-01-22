package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
	"github.com/wung-s/gotv/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// Auth0APIAudience identifies the server in Auth0
var Auth0APIAudience = []string{"https://gotv.com"}
var TwilioAccountSid = "AC23a198769cd7d761edb60783eccfa4c2"
var TwilioAuthToken = "f57fbe7303e5d9436e36dc4eafbf6796"
var TwilioNumber = "+15139121062"

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

		c := cors.AllowAll()

		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_gotv_session",
			PreWares:     []buffalo.PreWare{c.Handler},
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		app.Use(Authenticate)
		app.Middleware.Skip(Authenticate, HomeHandler)
		app.Middleware.Skip(Authenticate, RecruitersMembersSearch)

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(CurrentUserSetter)
		app.Middleware.Skip(CurrentUserSetter, HomeHandler)
		app.Middleware.Skip(CurrentUserSetter, RecruitersMembersSearch)

		app.GET("/", HomeHandler)

		// app.Resource("/members", MembersResource{})
		app.GET("/members", MembersList)
		app.PUT("/members/{id}", MembersUpdate)
		app.POST("/members/upload", MembersUpload)
		app.GET("/members/search", MembersSearch)

		app.POST("/users", UsersCreate)
		app.GET("/users", UsersList)
		app.PUT("/users/{id}", UsersUpdate)
		app.GET("/users/current", UsersCurrent)

		app.GET("/roles", RolesList)

		app.GET("/polls", PollsList)
		app.GET("/recruiters", RecruitersList)
		app.PUT("/recruiters/{id}", RecruitersUpdate)
		app.GET("/recruiters/{id}", RecruitersShow)
		app.GET("/recruiters/{id}/members/search", RecruitersMembersSearch)
	}

	return app
}
