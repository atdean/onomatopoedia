package webserver

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// App is the central dependency for all route controllers and config
type App struct {
	SqlPool         *sqlx.DB
	RedisConn 		redis.Conn
	Router			http.Handler
	IndexController *IndexController
	AuthController 	*AuthController
	EntryController *EntryController
}

func InitApp(sqlPool *sqlx.DB, redisConn redis.Conn) *App {
	app := &App{
		SqlPool: sqlPool,
		RedisConn: redisConn,
	}

	app.IndexController = NewIndexController(app)
	app.AuthController = NewAuthController(app)
	app.EntryController = newEntryController(app)

	app.Router = app.initRoutes()

	return app
}

func (app *App) initRoutes() http.Handler {
	router := httprouter.New()

	router.GET("/", app.IndexController.GetIndexHandler)

	router.GET("/entries/:slug", app.EntryController.GetSingleEntryHandler)

	router.GET("/login", app.AuthController.GetLoginHandler)
	router.POST("/login", app.AuthController.PostLoginHandler)

	router.ServeFiles("/static/*filepath", http.Dir("./public"))

	return router
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

