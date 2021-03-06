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
	"school-web-app/student"
	"school-web-app/user"
)

func main() {
	var config struct {
		Port int `default:"9000"`
	}
	envconfig.MustProcess("BACKEND", &config)

	db := mongo.GetDB()
	userDao := user.NewDao(db)
	authSvc := auth.NewService(db, userDao)
	studentSvc := student.NewService(db)
	routers := []server.Router{
		&rest.Users{Dao: userDao},
		&rest.Auth{Service: authSvc},
		&rest.Students{Service: studentSvc},
	}
	address := fmt.Sprintf(":%d", config.Port)
	logrus.WithField("port", config.Port).Info("starting HTTP server")
	handler := server.New(routers, authSvc)
	logrus.Fatal(http.ListenAndServe(address, handler))
}
