package student

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "student"

type Dao struct {
	db *mongo.Database
}

func NewDao(db *mongo.Database) *Dao {
	return &Dao{db}
}

func (d *Dao) Get(teacherID string) ([]Student, error) {
	c := d.db.Collection(collection)
	filter := bson.M{"$or": []bson.M{{"teacherId": teacherID}, {"teacherId": bson.M{"$exists": 0}}}}
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	s := make([]Student, 0)
	err = cur.All(context.Background(), &s)
	return s, err
}
