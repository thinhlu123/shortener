package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestClient_Connect(t *testing.T) {
	var a Collection

	var client Client
	if err := client.Connect(Configuration{
		Address:  []string{"localhost"},
		Username: "123abc123",
		Password: "123abc123",
		AuthDB:   "admin",
		Ssl:      false,
	}); err != nil {
		t.Error(err)
	}

	defer client.Disconnect()

	dbTest := client.InitDatabase("test", func(db *mongo.Database) error {
		a.Init("test_col", db)

		return nil
	})
	dbTest.InitCollection()

	doc := bson.M{"a": "b"}
	_, err := a.GetCollection().InsertOne(context.TODO(), doc)
	if err != nil {
		t.Error(err)
	}

	qRs := a.GetCollection().FindOne(context.TODO(), doc)
	if qRs.Err() != nil {
		t.Error(qRs)
	}
}
