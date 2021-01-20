package student

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "student"

type dao struct {
	db *mongo.Database
}

func newDao(db *mongo.Database) *dao {
	return &dao{db}
}

func (d *dao) Get(teacherID string, ids []primitive.ObjectID) ([]Student, error) {
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

func (d *dao) DeleteMany(teacherID string, ids []primitive.ObjectID) error {
	c := d.db.Collection(collection)
	filter := bson.M{"teacherId": teacherID, "_id": bson.M{"$in": ids}}
	res, err := c.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	if int(res.DeletedCount) != len(ids) {
		logrus.WithFields(logrus.Fields{
			"deleted":   res.DeletedCount,
			"requested": len(ids),
		}).Warn("failed to delete some of the students")
	}
	return nil
}

func (d *dao) CreateMany(students []Student) ([]Student, error) {
	c := d.db.Collection(collection)
	documents := make([]interface{}, len(students))
	for i := range students {
		students[i].ID = primitive.NewObjectID()
		documents[i] = students[i]
	}
	res, err := c.InsertMany(context.Background(), documents)
	if err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, len(res.InsertedIDs))
	for _, id := range res.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID))
	}
	return d.Get(students[0].TeacherID, ids)
}

func (d *dao) Update(teacherID string, student *Student) (*Student, error) {
	c := d.db.Collection(collection)
	filter := bson.M{"teacherId": teacherID, "_id": student.ID}
	update := bson.M{"$set": bson.M{
		"name":    student.Name,
		"surname": student.Surname,
		"age":     student.Age,
	}}
	_, err := c.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	res, err := d.Get(teacherID, []primitive.ObjectID{student.ID})
	if err != nil {
		return nil, err
	}
	return &res[0], err
}
