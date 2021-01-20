package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"

	"school-web-app/auth"
)

type TokenVerifier interface {
	Verify(token string) (*auth.Claims, error)
}

func New(routers []Router, verifier TokenVerifier) *Server {
	return &Server{newRouter(routers, verifier)}
}

type Server struct {
	router *httprouter.Router
}

func newRouter(routers []Router, verifier TokenVerifier) *httprouter.Router {
	muxRouter := httprouter.New()
	muxRouter.MethodNotAllowed = http.HandlerFunc(methodNotAllowedHandler)
	muxRouter.NotFound = http.HandlerFunc(notFoundHandler)
	chain := alice.New(timeLoggerHandler, basicHandler)
	for _, router := range routers {
		for _, route := range router.Routes() {
			ch := chain
			ch = ch.Append(route.PreHandlers...)
			if route.RequiresAuthorization {
				ch = ch.Append(authorizationHandler(verifier))
			}
			handler := ch.Then(route.HandlerFunc)
			muxRouter.Handler(route.Method, route.Pattern, handler)
		}
	}
	return muxRouter
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

func authorizationHandler(verifier TokenVerifier) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				logrus.WithField("token", token).Error("invalid token provided")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = strings.TrimPrefix(token, "Bearer ")
			claims, err := verifier.Verify(token)
			if err != nil {
				logrus.WithField("token", token).WithError(err).Error("invalid token provided")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			logrus.WithField("expiresAt", time.Unix(claims.ExpiresAt, 0).Format(time.RFC3339)).Info("authorized token")
			ctx := context.WithValue(r.Context(), "userId", claims.UserID)

			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
