package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
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

type PremiumUser struct {
	UserId       string             `bson:"user_id"`
	PremiumLevel int64              `bson:"premium_level"`
	MemberSince  primitive.DateTime `bson:"member_since"`
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
