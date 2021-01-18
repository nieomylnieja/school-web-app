package server

import (
	"net/http"

	"github.com/justinas/alice"
)

type Router interface {
	Routes() []Route
}

type Route struct {
	Name                  string
	Method                string
	Pattern               string
	PreHandlers           []alice.Constructor
	HandlerFunc           http.HandlerFunc
	RequiresAuthorization bool
}
