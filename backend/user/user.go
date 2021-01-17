package user

import (
	"errors"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	BirthDate string             `json:"birthDate" bson:"birthDate"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (u User) Validate() error {
	errs := strings.Builder{}
	if (len(u.Email) < 3 && len(u.Email) > 254) || !emailRegex.MatchString(u.Email) {
		errs.WriteString("invalid email provided\n")
	}
	required := func(s, field string) {
		if s == "" {
			errs.WriteString(field)
			errs.WriteString(" is required\n")
		}
	}
	required(u.Name, "name")
	required(u.Surname, "surname")
	required(u.BirthDate, "birth date")
	if errs.Len() != 0 {
		return errors.New(errs.String())
	}
	return nil
}
