package main

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	err := app.readJSON(w, r, requestPayload)
	if err != nil {
		log.Panic(err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    EmailContact{Email: requestPayload.From},
		To:      EmailContact{Email: requestPayload.To},
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Panic(err)
		app.errorJSON(w, err)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Sent to %s", requestPayload.To),
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)

}
