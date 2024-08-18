package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

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

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json to sent to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(http.MethodPost, "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// define custom net transport options
	var customTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 2 * time.Second,
	}

	// do not use default http client
	client := &http.Client{
		Timeout:   time.Second * 5,
		Transport: customTransport,
	}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// The client must close the response body when finished with it
	defer response.Body.Close()

	// verify status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	var jsonFromService jsonResponse

	// decode
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send error json back
	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload = jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
