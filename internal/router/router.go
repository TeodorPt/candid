// Copyright 2014 Canonical Ltd.

// The router package implements an HTTP request router for the id service.
package router

import (
	"net/http"
	"strings"

	"gopkg.in/errgo.v1"
)

// New returns an identity service router that will route requests to
// the given handlers.
//
// The handlers argument maps paths to handlers used to handle the paths.
// Path matching is by matched longest-prefix (the same as http.ServeMux).
// Note that, unlike http.ServeMux, the prefix is stripped from the URL path
// before the hander is invoked, matching the behaviour of the other handlers.
func New(handlers map[string]http.Handler) *Router {
	r := &Router{
		handlers: handlers,
	}
	mux := NewServeMux()
	for path, handler := range r.handlers {
		path = "/" + path
		prefix := strings.TrimSuffix(path, "/")
		mux.Handle(path, http.StripPrefix(prefix, handler))
	}
	r.mux = mux
	return r
}

// Router represents an identity service HTTP request router.
type Router struct {
	handlers map[string]http.Handler
	mux      http.Handler
}

// ServeHTTP implements http.Handler.ServeHTTP.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		WriteError(w, errgo.Notef(err, "cannot parse form"))
		return
	}
	r.mux.ServeHTTP(w, req)
}

// Handlers returns the set of handlers that the router was created with.
// This should not be changed.
func (r *Router) Handlers() map[string]http.Handler {
	return r.handlers
}