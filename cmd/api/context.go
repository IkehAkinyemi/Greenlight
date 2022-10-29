package main

import (
	"context"
	"net/http"

	"github.com/lighten/internal/data"
)

type contextKey string

var usercontextKey =  contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), usercontextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(usercontextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}