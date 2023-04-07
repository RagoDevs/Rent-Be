package main

import (
	"context"
	"net/http"

	"hmgt.hopertz.me/internal/data"
)

type contextKey string

const adminContextKey = contextKey("admin")

func (app *application) contextSetAdmin(r *http.Request, admin *data.Admin) *http.Request {
	ctx := context.WithValue(r.Context(), adminContextKey, admin)
	return r.WithContext(ctx)
}

func (app *application) contextGetAdmin(r *http.Request) *data.Admin {
	admin, ok := r.Context().Value(adminContextKey).(*data.Admin)
	if !ok {
		panic("missing admin value in request context")
	}
	return admin
}
