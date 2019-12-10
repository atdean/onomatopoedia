package webserver

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

// App is the central dependency for all route controllers and config
type App struct {
	SqlPool         *sqlx.DB
	RedisConn 		*redis.Conn
	IndexController *IndexController
	AuthController 	*AuthController
}

func InitApp(sqlPool *sqlx.DB, redisConn *redis.Conn) *App {
	app := &App{
		SqlPool: sqlPool,
		RedisConn: redisConn,
	}

	app.IndexController = NewIndexController(app)
	app.AuthController = NewAuthController(app)

	return app
}