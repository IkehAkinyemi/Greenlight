package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Greenlight/internal/data"
)

// showMovie maps to the "GET /v1/movies/:id" endpoint.
func (app *application) showMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.retrieveIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
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
		app.serverErrorResponse(w, r, err)
	}
}

// createMovie maps to the "POST /v1/movies" endpoint.
func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	var movie struct{
		Title string `json:"title"`
		Runtime int32 `json:"runtime"`
		Year int32 `json:"year"`
		Genre []string `json:"genre"`
	}

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v", movie)
}