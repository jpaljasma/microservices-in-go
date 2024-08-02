package main

import (
	"net/http"
	"time"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:     false,
		Message:   "Hit the broker",
		Timestamp: time.Now().UTC().String(),
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) Index(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:     false,
		Message:   "Welcome to broker api homepage",
		Timestamp: time.UTC.String(),
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
