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
		Keys: bson.D{{
			Key:   "u",
			Value: 1,
		}, {
			Key:   "_id",
			Value: -1,
		}},
	}, {
		Keys: bson.M{
			"e": 1,
		},
	}, {
		Keys: bson.M{
			"ca": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32((3 * 24 * time.Hour).Seconds())),
	}})
	if err != nil {
		return
	}

	_, err = db.Collection(CollCachedExport).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"m": 1,
		},
		Options: options.Index().SetUnique(true),
	}, {
		Keys: bson.M{
			"ca": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32((3 * 24 * time.Hour).Seconds())),
	}})
	if err != nil {
		return
	}

	_, err = db.Collection(CollProcessingExport).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"l": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32(time.Minute.Seconds())),
	}, {
		Keys: bson.M{
			"e": 1,
		},
		Options: options.Index().SetUnique(true),
	}})
	if err != nil {
		return
	}

	_, err = db.Collection(CollRow).Indexes().CreateMany(ctx, []mongo.IndexModel{{
		Keys: bson.M{
			"e": 1,
		},
	}, {
		Keys: bson.M{
			"ca": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(int32((10 * time.Hour).Seconds())),
	}})
	return
}
