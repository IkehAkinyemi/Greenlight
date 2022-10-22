package main

import (
	"net/http"
)

// healthcheck maps to "GET /v1/healthcheck". Return info about the server state.
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"enviroment": app.config.env,
			"version":    version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encounter an issue and can't process request", http.StatusInternalServerError)
	}
}
