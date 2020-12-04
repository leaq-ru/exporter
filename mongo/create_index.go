package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func createIndex(db *mongo.Database) (err error) {
	ctx := context.Background()

	_, err = db.Collection(CollFile).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"u": 1,
		},
		Options: options.Index().SetUnique(true),
	}, {
		Keys: bson.M{
			"c": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32((24 * time.Hour).Seconds())),
	}})
	return
}
