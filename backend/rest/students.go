package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"school-web-app/server"
	"school-web-app/student"
)

type Students struct {
	Service *student.Service
}

func (s *Students) Routes() []server.Route {
	return []server.Route{
		{
			Method:                "GET",
			Pattern:               "/students",
			HandlerFunc:           s.Get,
			RequiresAuthorization: true,
		},
		{
			Method:                "PUT",
			Pattern:               "/students",
			HandlerFunc:           s.Put,
			RequiresAuthorization: true,
		},
	}
}

func (s *Students) Get(w http.ResponseWriter, r *http.Request) {
	teacherID := r.Context().Value("userId").(string)
	res, err := s.Service.Get(teacherID)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch students")
		respond(w, "internal server error", http.StatusInternalServerError)
		return
	}
	respond(w, res, http.StatusOK)
}

func (s *Students) Put(w http.ResponseWriter, r *http.Request) {
	teacherID := r.Context().Value("userId").(string)
	var payload []student.Student
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.WithError(err).Warn("failed to decode json body")
		respond(w, nil, http.StatusBadRequest)
		return
	}
	res, err := s.Service.Put(teacherID, payload)
	if err != nil {
		respondWithError(w, err)
		return
	}
	respond(w, res, http.StatusCreated)
}
