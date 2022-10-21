package main

import (
	"fmt"
	"net/http"
)

// showMovie maps to the "GET /v1/movies/:id" endpoint.
func (app *application) showMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveIDParam(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "showing details for movie-id %d\n", id)
}

// createMovie maps to the "POST /v1/movies" endpoint.
func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new movie")
}