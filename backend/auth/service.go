package auth

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"school-web-app/errz"
	"school-web-app/user"
)

type userGetter interface {
	GetByEmail(email string) (*user.User, error)
	GetByID(id primitive.ObjectID) (*user.User, error)
}

type Service struct {
	dao        *dao
	userGetter userGetter
	ts         *tokenStore
}

func NewService(db *mongo.Database, userGetter userGetter) *Service {
	return &Service{
		dao:        newDao(db),
		userGetter: userGetter,
		ts:         newTokenStore(),
	}
}

func (s *Service) SavePassword(password *Password) (*Password, error) {
	p, err := s.dao.SavePassword(password)
	return p, err
}

func (s *Service) SignIn(payload *SignInPayload) (*Token, error) {
	u, err := s.userGetter.GetByEmail(payload.Email)
	if err != nil {
		return nil, err
	}
	p, err := s.dao.GetPasswordByUserID(u.ID)
	if err != nil {
		return nil, err
	}
	if p.Password != payload.Password {
		return nil, errz.New(errz.Unauthorized, "invalid password", nil)
	}
	return s.ts.Generate(u.ID.Hex()), nil
}

func (s *Service) Verify(token string) (string, error) {
	claims, err := s.ts.Validate(token)
	if err != nil {
		return "", err
	}
	id, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return "", errz.New(errz.Unauthorized, "invalid token provided", err)
	}
	_, err = s.userGetter.GetByID(id)
	if err != nil {
		if e, ok := err.(*errz.Error); ok && e.Type() == errz.NotFound {
			return "", errz.New(errz.Unauthorized, "user provided in token doesn't exist", err)
		}
		return "", err
	}
	return id.Hex(), nil
}
