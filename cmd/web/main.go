package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/webserver"
)

func main() {
	dbConn, err := sql.Open("mysql", "onomatopoedia:onomatopoedia@/onomatopoedia")
	if err != nil {
		log.Fatalln(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	app := webserver.App{
		SqlPool: dbConn,
		IndexController: webserver.NewIndexController(dbConn),
	}

	fmt.Println("HTTP server started... listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", initRoutes(&app)))
}

func initRoutes(app *webserver.App) http.Handler {
	router := httprouter.New()

	router.GET("/", app.IndexController.GetIndexHandler)

	router.ServeFiles("/static/*filepath", http.Dir("./public"))

	return router
}
