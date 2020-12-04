package file

import "go.mongodb.org/mongo-driver/mongo"

type Model struct {
	db    *mongo.Database
	files *mongo.Collection
}
