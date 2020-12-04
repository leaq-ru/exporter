package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StartSession func(...*options.SessionOptions) (mongo.Session, error)
