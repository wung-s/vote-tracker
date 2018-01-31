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

// Response defines the response message structure
type Response struct {
	Message string `json:"message"`
}

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		if ENV == "development" {
			envy.Load()
		}

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

		InitializeFirebase()
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(Authenticate)
		app.Middleware.Skip(Authenticate, HomeHandler)
		app.Middleware.Skip(Authenticate, RecruitersMembersSearch)
		app.Middleware.Skip(Authenticate, RecruitersShow)

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.

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
		app.GET("/polls/streets", PollsStreets)
		app.GET("/polls/{id}", PollsShow)

		app.GET("/recruiters", RecruitersList)
		app.PUT("/recruiters/{id}", RecruitersUpdate)
		app.PUT("/recruiters/{id}/invite", RecruitersInvite)
		app.POST("/recruiters/inviteAll", RecruitersInviteAll)
		app.GET("/recruiters/{id}", RecruitersShow)
		app.GET("/recruiters/{id}/members/search", RecruitersMembersSearch)

		app.POST("/members/{id}/dispositions", DispositionsCreate)
	}

	return app
}
