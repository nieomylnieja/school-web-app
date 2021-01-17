package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"school-web-app/server"
	"school-web-app/user"
)

type Users struct {
	Dao *user.Dao
}

func (u *Users) Routes() []server.Route {
	return []server.Route{
		{
			Method:      "POST",
			Pattern:     "/users",
			HandlerFunc: u.Create,
		},
	}
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var payload user.User
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
	res, err := u.Dao.Create(&payload)
	if err != nil {
		logrus.WithError(err).Error("failed to create user")
		respond(w, "internal server error", http.StatusInternalServerError)
		return
	}
	respond(w, res, http.StatusCreated)
}
