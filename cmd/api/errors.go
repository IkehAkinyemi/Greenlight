package main

import (
	"fmt"
	"net/http"
)

// The logError method is a generic helper for logging an error message and
// additional information from the request including the HTTP method and URL.
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// serverErrorResponse() method reports runtime errors/problems.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	msg := "the server encountered an error and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
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

// failedValidationResponse reports errors from JSON validation
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// editConflictResponse reports edit conflict, like data race.
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, msg)

}

// rateLimitExceededResponse reports rate limiting errors.
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	msg := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, msg)
}

// invalidCredentialResponse reports user authentication errors.
func (app *application) invalidCredentialResponse(w http.ResponseWriter, r *http.Request) {
	msg := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

// invalidAuthenticationTokenResponse reports user authentication errors in regards to token
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	// Keeps a reminder for the client about the bearer token
	w.Header().Add("WWW-Authentication", "Bearer")
	msg := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

// authenticationRequiredResponse reports error relating to token-based authentation.
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	msg := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, msg)
}

// inactiveAccountResponse reports error if user is not activated yet.
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	msg := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, msg)
}

// notPermittedResponse reports error if user isn't authorized for a particular resource.
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	msg := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, msg)
}
