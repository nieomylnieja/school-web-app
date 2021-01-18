package rest

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"school-web-app/server"
	"school-web-app/student"
)

type Students struct {
	Dao *student.Dao
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
			HandlerFunc:           s.Upsert,
			RequiresAuthorization: true,
		},
	}
}

func (s *Students) Get(w http.ResponseWriter, r *http.Request) {
	teacherID := r.Context().Value("userId").(string)
	res, err := s.Dao.Get(teacherID)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch students")
		respond(w, "internal server error", http.StatusInternalServerError)
		return
	}
	logrus.Info(res)
	respond(w, res, http.StatusOK)
}

func (s *Students) Upsert(w http.ResponseWriter, r *http.Request) {
	// TODO how should it be handled? That's a question to frontend really
}
