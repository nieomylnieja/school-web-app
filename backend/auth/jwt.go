package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"school-web-app/errz"
)

type Token struct {
	Token string `json:"token"`
}

type tokenStore struct {
	SecretKey string        `split_words:"true" required:"true"`
	Expiry    time.Duration `default:"15m"`
}

func newTokenStore() *tokenStore {
	t := tokenStore{}
	envconfig.MustProcess("jwt_token", &t)
	return &t
}

type claims struct {
	Authorized bool   `json:"authorized"`
	UserID     string `json:"userId"`
	jwt.StandardClaims
}

func (t *tokenStore) Generate(userID string) *Token {
	claims := claims{
		Authorized:     true,
		UserID:         userID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: int64(t.Expiry)},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(t.SecretKey))
	if err != nil {
		logrus.WithField("userID", userID).WithError(err).Panic("failed to sign the token")
	}
	return &Token{token}
}

func (t *tokenStore) Validate(tokenString string) (*claims, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errz.New(errz.Invalid, "invalid token provided", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, errz.New(errz.Unauthorized, "expired token", err)
			}
		}
		return nil, err
	}
	return claims, nil
}
