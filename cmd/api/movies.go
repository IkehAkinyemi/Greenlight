package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Greenlight/internal/data"
)

// showMovie maps to the "GET /v1/movies/:id" endpoint.
func (app *application) showMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveIDParam(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID: id,
		Title: "Halloween",
		Runtime: 120,
		CreatedAt: time.Now(),
		Version: 1,
		Genre: []string{"rom-com", "horror", "sci-fi"},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// createMovie maps to the "POST /v1/movies" endpoint.
func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new movie")
}