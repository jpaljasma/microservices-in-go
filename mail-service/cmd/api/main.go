package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "80"

func main() {

	app := Config{}

	log.Printf("Starting mail service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start webserver
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
