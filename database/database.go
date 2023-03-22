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
	UserId int64              `bson:"user_id"`
	Name   string
	Bees   map[int]*hive.Bee
}

type TriviaQuestion struct {
	ID         primitive.ObjectID `bson:"_id"`
	Difficulty int
	Question   string
	Answer     string
	Incorrect  []string
}

type LeaderboardUser struct {
	UserId        int64 `bson:"user_id"`
	Username      string
	Discriminator string
	TriviaPoints  int   `bson:"trivia_points"`
	BubbleTime    int64 `bson:"bubble_time"`
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
