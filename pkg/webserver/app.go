package webserver

import "github.com/jmoiron/sqlx"

// App is the central dependency for all route controllers and config
type App struct {
	SqlPool         *sqlx.DB
	IndexController *IndexController
}
