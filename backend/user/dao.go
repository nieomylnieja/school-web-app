package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"school-web-app/errz"
)

const collection = "user"

type Dao struct {
	db *mongo.Database
}

func NewDao(db *mongo.Database) *Dao {
	return &Dao{db}
}

func (d *Dao) Create(user *User) (*User, error) {
	c := d.db.Collection(collection)
	user.ID = primitive.NewObjectID()
	if _, err := c.InsertOne(context.Background(), user); err != nil {
		return nil, err
	}
	return d.GetByID(user.ID)
}

func (d *Dao) GetByID(id primitive.ObjectID) (*User, error) {
	return d.getByFilter(bson.M{"_id": id})
}

func (d *Dao) GetByEmail(email string) (*User, error) {
	return d.getByFilter(bson.M{"email": email})
}

func (d *Dao) getByFilter(filter bson.M) (*User, error) {
	c := d.db.Collection(collection)
	var u User
	err := c.FindOne(context.Background(), filter).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, errz.New(errz.NotFound, "user doesn't exist", err)
	}
	return &u, err
}
