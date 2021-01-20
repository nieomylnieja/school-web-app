package student

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Age       int                `json:"age" bson:"age"`
	TeacherID string             `json:"teacherId,omitempty" bson:"teacherId,omitempty"`
}

func (s Student) Validate() error {
	errs := strings.Builder{}
	required := func(s interface{}, field string) {
		if s == "" || s == 0 {
			errs.WriteString(field)
			errs.WriteString(" is required\n")
		}
	}
	required(s.Name, "name")
	required(s.Surname, "surname")
	required(s.Age, "age")
	if errs.Len() != 0 {
		return errors.New(errs.String())
	}
	return nil
}
