package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// A standard json response type used in our broker service
type jsonResponse struct {
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Timestamp string `json:"ts"`
	Data      any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	// ensure json payload length is less than 1 megabyte
	const maxBytes = 1_048_576                               // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // truncate

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// ensure there's only a single json valye
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)

	if err != nil {
		return err
	}

	return nil
}

// Outputs the error into http response in standard format using [jsonResponse]
// See also [writeJSON]
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Timestamp = time.UTC.String()
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
