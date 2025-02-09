package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}

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

func createMail() Mail {

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:     os.Getenv("MAIL_DOMAIN"),
		Host:       os.Getenv("MAIL_HOST"),
		Port:       port,
		Username:   os.Getenv("MAIL_USERNAME"),
		Password:   os.Getenv("MAIL_PASSWORD"),
		Encryption: os.Getenv("MAIL_ENCRYPTION"),
		From: EmailContact{
			Name:  os.Getenv("MAIL_FROM_NAME"),
			Email: os.Getenv("MAIL_FROM_EMAIL"),
		},
	}

	return m
}
