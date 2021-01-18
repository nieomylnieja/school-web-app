package auth

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "password"

type dao struct {
	cipher *cipher
	db     *mongo.Database
}

func newDao(db *mongo.Database) *dao {
	return &dao{
		db:     db,
		cipher: newCipher(),
	}
}

func (d *dao) SavePassword(password *Password) (*Password, error) {
	c := d.db.Collection(collection)
	password.ID = primitive.NewObjectID()
	password.Password = d.cipher.Encrypt(password.Password)
	if _, err := c.InsertOne(context.Background(), password); err != nil {
		return nil, err
	}
	return d.GetPasswordByUserID(password.UserID)
}

func (d *dao) GetPasswordByUserID(userID primitive.ObjectID) (*Password, error) {
	c := d.db.Collection(collection)
	var p Password
	if err := c.FindOne(context.Background(), bson.M{"userId": userID}).Decode(&p); err != nil {
		return nil, err
	}
	p.Password = d.cipher.Decrypt(p.Password)
	return &p, nil
}
