package rest

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"school-web-app/errz"
)

func respondWithError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	statusText := []byte("internal server error")
	e, ok := err.(*errz.Error)
	if ok {
		statusCode = e.HttpStatusCode()
		statusText = []byte(e.HttpStatusText())
		w.WriteHeader(statusCode)
		if _, err = w.Write(statusText); err != nil {
			logrus.WithError(err).Error("failed to write response body")
		}
	}
	logrus.WithField("statusCode", statusCode).WithError(err).Error(string(statusText))
}

// respond trusts that the programmer will only try to encode either strings or objects
func respond(w http.ResponseWriter, v interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	if v == nil {
		return
	}
	write := func(data []byte) {
		if _, err := w.Write(data); err != nil {
			logrus.WithError(err).Error("failed to write response body")
		}
	}
	switch v.(type) {
	case string:
		write([]byte(v.(string)))
	default:
		data, err := json.Marshal(v)
		if err != nil {
			logrus.WithError(err).Error("failed to marshall response body")
		}
		write(data)
	}
}
