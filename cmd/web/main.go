package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"

	"github.com/atdean/onomatopoedia/pkg/webserver"
)

func main() {
	sqlPool, err := sqlx.Connect("mysql", "onomatopoedia:onomatopoedia@/onomatopoedia")
	if err != nil {
		log.Fatalln(err)
	}

	redisConn, err := redis.DialURL("redis://localhost")
	if err != nil {
		log.Fatalln(err);
	}

	//app := webserver.App{
	//	SqlPool: sqlPool,
	//	RedisConn: &redisConn,
	//	IndexController: webserver.NewIndexController(sqlPool),
	//	AuthController: webserver.NewAuthController(sqlPool, &redisConn),
	//}

	app := webserver.InitApp(sqlPool, &redisConn)

	fmt.Println("HTTP server started... listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", initRoutes(app)))
}

func initRoutes(app *webserver.App) http.Handler {
	router := httprouter.New()

	router.GET("/", app.IndexController.GetIndexHandler)

	router.GET("/login", app.AuthController.GetLoginHandler)
	router.POST("/login", app.AuthController.PostLoginHandler)

	router.ServeFiles("/static/*filepath", http.Dir("./public"))

	return router
}
