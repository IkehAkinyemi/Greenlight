package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// envelope wraps the JSON response.
type envelope map[string]interface{}

// retrieveIDParam returns the "id" URL parameter from the current request context,
// then convert it to an integer and return it.
func (app *application) retrieveIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// writeJSON send responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data envelope, header http.Header) error {
	resp, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	resp = append(resp, '\n')
	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(resp)

	return nil
}