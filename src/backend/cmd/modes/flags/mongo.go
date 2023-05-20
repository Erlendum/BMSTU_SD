package flags

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoFlags struct {
	Uri    string `mapstructure:"uri"`
	DBName string `mapstructure:"dbname"`
}

func (m *MongoFlags) InitDB() (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(m.Uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database(m.DBName), nil
}
