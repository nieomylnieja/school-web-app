package auth

import (
	"errors"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Password struct {
	ID       primitive.ObjectID `json:"-" bson:"_id"`
	UserID   primitive.ObjectID `json:"userId" bson:"userId"`
	Password string             `json:"password" bson:"password"`
}

var passwordRegex = regexp.MustCompile(`\d`)

func (p Password) Validate() error {
	if len(p.Password) < 8 || !passwordRegex.MatchString(p.Password) {
		return errors.New("password must contain at least one number and be longer than 7 characters")
	}
	return nil
}
