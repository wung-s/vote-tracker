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
			// Uncomment and fix the issue with Redis URL for Redis-based background job
			// Worker: gwa.New(gwa.Options{
			// 	Pool: &redis.Pool{
			// 		MaxActive: 5,
			// 		MaxIdle:   5,
			// 		Wait:      true,
			// 		Dial: func() (redis.Conn, error) {
			// 			return redis.Dial("tcp", os.Getenv("REDIS_URL"))
			// 		},
			// 	},
			// 	Name:           "gotv",
			// 	MaxConcurrency: 25,
			// }),
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		InitializeGoogleMaps()

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		tx := middleware.PopTransaction(models.DB)

		app.Use(tx)
		// app.Use(Authenticate)
		app.Use(RestrictedHandlerMiddleware)

		app.Middleware.Skip(RestrictedHandlerMiddleware, HomeHandler)
		app.Middleware.Skip(RestrictedHandlerMiddleware, LoginHandler)
		app.Middleware.Skip(RestrictedHandlerMiddleware, RecruitersMembersSearch)
		app.Middleware.Skip(RestrictedHandlerMiddleware, RecruitersShow)
		app.Middleware.Skip(tx, MembersUpload)
		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		app.GET("/", HomeHandler)
		app.POST("/login", LoginHandler)

		app.GET("/members", MembersList)
		app.PUT("/members/{id}", MembersUpdate)
		app.POST("/members/upload", MembersUpload)
		app.GET("/members/search", MembersSearch)

		app.POST("/users", UsersCreate)
		app.GET("/users", UsersList)
		app.PUT("/users/{id}", UsersUpdate)
		app.GET("/users/current", UsersCurrent)
		app.GET("/users/{user_id}", UsersShow)
		// app.POST("/signup", UsersSignUp)

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

		app.GET("/ride_requests", RideRequestList)
		app.POST("/members/{id}/dispositions", DispositionsCreate)
		app.POST("/members/{member_id}/ride_requests", RideRequestsCreate)
		app.GET("/members/{member_id}/ride_requests", MemberRideRequestsShow)
		app.PUT("/ride_requests/{id}", RideRequestUpdate)

	}

	return app
}
