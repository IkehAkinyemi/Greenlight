package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Greenlight/internal/data"
	"github.com/Greenlight/internal/validator"
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
		Genres: []string{"rom-com", "horror", "sci-fi"},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createMovie maps to the "POST /v1/movies" endpoint.
func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	var input struct{
		Title string `json:"title"`
		Runtime data.Runtime `json:"runtime"`
		Year int32 `json:"year"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title: input.Title,
		Runtime: input.Runtime,
		Year: input.Year,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v", movie)
}

