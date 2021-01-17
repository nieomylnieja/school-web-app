package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/kelseyhightower/envconfig"
)

var config mongoConfig

type mongoConfig struct {
	Host            string        `default:"127.0.0.1"`
	Port            int           `default:"27017"`
	MaxPoolSize     uint64        `split_words:"true" default:"0"`
	MaxConnIdleTime time.Duration `split_words:"true" default:"30m"`
	// credentials
	Database string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

func (c mongoConfig) getMongoURL() string {
	u := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
	logrus.Info(u)
	return u
}

func GetDB() *mongo.Database {
	ctx := context.Background()
	envconfig.MustProcess("mongo", &config)
	mongoOpts := options.Client().
		SetMaxConnIdleTime(config.MaxConnIdleTime).
		SetMaxPoolSize(config.MaxPoolSize).
		SetReadPreference(readpref.Primary()).
		ApplyURI(config.getMongoURL())
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		logrus.WithError(err).Panic("Error occurred while connecting to mongoDB")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logrus.WithError(err).Panic("Error occurred while verifying mongoDB connection")
	}
	return client.Database(config.Database)
}
