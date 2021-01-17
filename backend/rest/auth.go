package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"school-web-app/auth"
	"school-web-app/server"
)

type Auth struct {
	Service *auth.Service
}

func (a *Auth) Routes() []server.Route {
	return []server.Route{
		{
			Method:      "POST",
			Pattern:     "/auth/signup",
			HandlerFunc: a.SignUp,
		},
		{
			Method:      "POST",
			Pattern:     "/auth/signin",
			HandlerFunc: a.SignIn,
		},
		{
			Method:      "GET",
			Pattern:     "/auth/verify",
			HandlerFunc: a.Verify,
		},
	}
}

func (a *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload auth.Password
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.WithError(err).Warn("failed to decode json body")
		respond(w, nil, http.StatusBadRequest)
		return
	}
	if err := payload.Validate(); err != nil {
		logrus.WithError(err).Warn("invalid body provided")
		respond(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := a.Service.SavePassword(&payload)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respond(w, res, http.StatusCreated)
}

func (a *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload auth.SignInPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.WithError(err).Warn("failed to decode json body")
		respond(w, nil, http.StatusBadRequest)
		return
	}
	res, err := a.Service.SignIn(&payload)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respond(w, res, http.StatusOK)
}

func (a *Auth) Verify(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := a.Service.Verify(token)
	if err != nil {
		respondWithError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
