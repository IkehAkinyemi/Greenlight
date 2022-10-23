package main

import (
	"fmt"
	"net/http"
)

// The logError method is a generic helper for logging an error message and
// additional information from the request including the HTTP method and URL.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// serverErrorResponse() method reports runtime errors/problems.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	msg := "the server encountered an error and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app * application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resourcec could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

// errorResponse method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, statusCode, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}