package main

import (
	"log"
)

type Config struct {
}

const webPort = "80"

func main() {

	msg := Message{}
	msg.From = EmailContact{Name: "John", Email: "test@test.com"}

	mail := &Mail{}
	err := mail.SendSMTPMessage(msg)
	if err != nil {
		log.Panic(err)
	}
	return

	// app := Config{}

	// log.Printf("Starting mail service on port %s\n", webPort)

	// // define http server
	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", webPort),
	// 	Handler: app.routes(),
	// }

	// // start webserver
	// err := srv.ListenAndServe()
	// if err != nil {
	// 	log.Panic(err)
	// }

}
