package student

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"school-web-app/errz"
)

type Service struct {
	dao *dao
}

func NewService(db *mongo.Database) *Service {
	return &Service{dao: newDao(db)}
}

func (s *Service) Get(teacherID string) ([]Student, error) {
	return s.dao.Get(teacherID, nil)
}

func (s *Service) Put(teacherID string, payload []Student) ([]Student, error) {
	errsBuf := strings.Builder{}
	for i := range payload {
		if err := payload[i].Validate(); err != nil {
			errsBuf.WriteString(err.Error() + "\n")
		}
	}
	if errsBuf.Len() != 0 {
		return nil, errz.New(errz.Invalid, errsBuf.String(), nil)
	}

	existing, err := s.dao.Get(teacherID, nil)
	if err != nil {
		return nil, err
	}
	if err = s.handlePutDeletion(teacherID, payload, existing); err != nil {
		return nil, err
	}
	return s.handlePutUpsert(teacherID, payload, existing)
}

func (s *Service) handlePutUpsert(teacherID string, payload, existing []Student) ([]Student, error) {
	var toBeUpdated, toBeCreated []Student
outerUpsert:
	for _, p := range payload {
		for _, e := range existing {
			if p.ID == e.ID {
				toBeUpdated = append(toBeUpdated, p)
				continue outerUpsert
			}
		}
		p.TeacherID = teacherID
		toBeCreated = append(toBeCreated, p)
	}
	var err error
	var created []Student
	if toBeCreated != nil {
		created, err = s.dao.CreateMany(toBeCreated)
		if err != nil {
			return nil, err
		}
	}
	updated := make([]Student, len(toBeUpdated))
	for _, u := range toBeUpdated {
		u, err := s.dao.Update(teacherID, &u)
		if err != nil {
			return nil, err
		}
		updated = append(updated, *u)
	}
	return append(created, updated...), nil
}

func (s *Service) handlePutDeletion(teacherID string, payload, existing []Student) error {
	var toBeDeleted []primitive.ObjectID
outerDelete:
	for _, e := range existing {
		for _, p := range payload {
			if e.ID == p.ID || p.ID == primitive.NilObjectID {
				continue outerDelete
			}
		}
		toBeDeleted = append(toBeDeleted, e.ID)
	}
	if toBeDeleted != nil {
		if err := s.dao.DeleteMany(teacherID, toBeDeleted); err != nil {
			return err
		}
	}
	return nil
}
