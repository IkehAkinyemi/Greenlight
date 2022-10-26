package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lighten/internal/validator"
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

// readJSON reads/parses request body. Also handles any possible error
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Restrict r.Body to 1MB
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)

	//Produces error if any unknown json fields is present
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)

	if err != nil {
		// types of expected errors
		var syntaxError *json.SyntaxError
		var unmarshaTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formatted JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formatted JSON")
		case errors.As(err, &unmarshaTypeError):
			if unmarshaTypeError.Field != "" {
				return fmt.Errorf("body contains badly-formatted JSON type for the field: %q", unmarshaTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formatted JSON type for the field: %d", unmarshaTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("the body contains unknown field %s", field)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		default:
			return err
		}
	}

	//second call to Decode to ensure the request body is just on
	// JSON value
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only a single JSON")
	}

	return nil
}

// readStr reads/parses the query string for a key's value
func (app *application) readStr(queryStr url.Values, key, defaultValue string) string {
	str := queryStr.Get(key)
	if str == "" {
		return defaultValue
	}
	return str
}

// readCSV parses the csv-like values provide in the query string
func (app *application) readCSV(queryStr url.Values, key string, defaultSlice []string) []string {
	csv := queryStr.Get(key)

	if csv == "" {
		return defaultSlice
	}
	return strings.Split(csv, ",")
}

// readInt parses integer values provided through the query string
func (app *application) readInt(queryStr url.Values, key string, defaultValue int, v *validator.Validator) int {
	str := queryStr.Get(key)
	if str == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(str)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return intValue
}

// The background() helper accepts an arbitrary function as a parameter.
func (app *application) backgroundJob(fn func()) {
	app.wg.Add(1)

	go func ()  {
		defer app.wg.Done()
		
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
