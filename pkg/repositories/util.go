package repositories

import (
	"io"
	"log"
)

func closeConnection(conn io.Closer) {
	if err := conn.Close(); err != nil {
		// Do we want to crash the application in this case?
		log.Fatalln(err)
	}
}