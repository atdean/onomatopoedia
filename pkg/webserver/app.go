package webserver

import "database/sql"

// App is the central dependency for all route controllers and config
type App struct {
	SqlPool         *sql.DB
	IndexController *IndexController
}
