package cached_export

import "go.mongodb.org/mongo-driver/mongo"

type Model struct {
	db            *mongo.Database
	cachedExports *mongo.Collection
}
