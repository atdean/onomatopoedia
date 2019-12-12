package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/atdean/onomatopoedia/pkg/webserver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
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


	app := webserver.InitApp(sqlPool, redisConn)

	fmt.Println("HTTP server started... listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", app))
}
