package main

import (
	"context"
	"net/http"

	"github.com/lighten/internal/data"
)

type contextKey string

var usercontextKey = contextKey("user")

// contextSetUser registers an authenticated user per connection
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), usercontextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves n authenticated user.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(usercontextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
