package common

import (
	"context"
	"errors"

	"alaninnovates.com/hive-bot/database"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PremiumLevel int

const (
	PremiumLevelFree PremiumLevel = iota
	PremiumLevelBuilder
)

func GetPremiumUsers(db database.Database, userID snowflake.ID) (database.PremiumUser, error) {
	var u database.PremiumUser
	err := db.Collection("premium-users").FindOne(context.Background(), bson.M{
		"user_id": userID.String(),
	}).Decode(&u)
	if err != nil {
		return database.PremiumUser{}, err
	}
	return u, nil
}

func GetPremiumLevel(db database.Database, userID snowflake.ID) (PremiumLevel, error) {
	u, err := GetPremiumUsers(db, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return PremiumLevelFree, nil
		}
		return -1, err
	}
	return PremiumLevel(u.PremiumLevel), nil
}
