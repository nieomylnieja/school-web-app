package server

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
)

func New(routers []Router) *Server {
	return &Server{
		router: newRouter(routers),
	}
}

type Server struct {
	router *httprouter.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,PATCH,PUT")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User, X-Tenant, X-User-Data, X-User-Location, X-Origin-System")
	}
	if r.Method == "OPTIONS" {
		return
	}
	s.router.ServeHTTP(w, r)
}

func newRouter(routers []Router) *httprouter.Router {
	muxRouter := httprouter.New()
	muxRouter.GlobalOPTIONS = http.HandlerFunc(globalOptionsHandler)
	muxRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowedHandler)
	muxRouter.NotFound = http.HandlerFunc(notFoundHandler)
	chain := alice.New(timeLoggerHandler, basicHandler)
	for _, router := range routers {
		for _, route := range router.Routes() {
			ch := chain
			ch = ch.Append(route.PreHandlers...)
			handler := ch.Then(route.HandlerFunc)
			muxRouter.Handler(route.Method, route.Pattern, handler)
		}
	}
	return muxRouter
}

func timeLoggerHandler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		inner.ServeHTTP(w, r)

		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"URI":    r.RequestURI,
		}).Infof("Request handled in: %v", time.Since(begin))
	})
}

func globalOptionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Access-Control-Request-Method") != "" {
		// Set CORS headers
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", r.Header.Get("Allow"))
		header.Set("Access-Control-Allow-Origin", "*")
	}
	w.WriteHeader(http.StatusNoContent)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
	return
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return
}

func basicHandler(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		inner.ServeHTTP(w, r)
	})
}
