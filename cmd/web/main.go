package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/webserver"
)

func main() {
	app := webserver.App{
		IndexController: webserver.IndexController{},
	}

	fmt.Println("HTTP server started... listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", initRoutes(&app)))
}

func initRoutes(app *webserver.App) http.Handler {
	router := httprouter.New()

	router.GET("/", app.IndexController.GetIndexHandler)

	return router
}
