package database

import (
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
}

type UserHive struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	Name   string
	Bees   map[int]*hive.Bee
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Collection(coll string) *mongo.Collection {
	return d.client.Database("hive-bot").Collection(coll)
}

func (d *Database) Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	d.client = client
	return client, nil
}
