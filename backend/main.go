package main

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"school-web-app/auth"
	"school-web-app/mongo"
	"school-web-app/rest"
	"school-web-app/server"
	"school-web-app/user"
)

func main() {
	var config struct {
		Port int `default:"9000"`
	}
	envconfig.MustProcess("BACKEND", &config)

	db := mongo.GetDB()
	userDao := user.NewDao(db)
	routers := []server.Router{
		&rest.Users{Dao: userDao},
		&rest.Auth{Service: auth.NewService(db, userDao)},
	}
	handler := server.New(routers)
	logrus.WithField("port", config.Port).Info("starting HTTP server")
	address := fmt.Sprintf("localhost:%d", config.Port)
	logrus.Fatal(http.ListenAndServe(address, handler))
}
