package webserver

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// App is the central dependency for all route controllers and config
type App struct {
	SqlPool         *sqlx.DB
	RedisConn 		redis.Conn
	Router			http.Handler
	SecretKey		string

	InfoLogger		*log.Logger
	ErrorLogger		*log.Logger

	IndexController *IndexController
	AuthController 	*AuthController
	EntryController *EntryController
}

func InitApp(sqlPool *sqlx.DB, redisConn redis.Conn) *App {
	app := &App{
		SqlPool: sqlPool,
		// TODO :: Make sure this isn't being copied and creating double connections
		RedisConn: redisConn,
		// TODO :: Pull this from an env var or flag, and make a CLI tool to regenerate it
		SecretKey: "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge",
	}

	app.InfoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLogger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.IndexController = NewIndexController(app)
	app.AuthController = NewAuthController(app)
	app.EntryController = NewEntryController(app)

	app.Router = app.initRoutes()

	return app
}

func (app *App) initRoutes() http.Handler {
	router := httprouter.New()

	router.GET("/", app.IndexController.GetIndexHandler)

	router.GET("/entries/:slug", app.EntryController.GetSingleEntryHandler)

	router.GET("/signup", app.AuthController.GetSignupHandler)
	router.POST("/signup", app.AuthController.PostSignupHandler)
	router.GET("/login", app.AuthController.GetLoginHandler)
	router.POST("/login", app.AuthController.PostLoginHandler)

	router.ServeFiles("/static/*filepath", http.Dir("./public"))

	return router
}


func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

