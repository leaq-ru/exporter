package processing_export

import "go.mongodb.org/mongo-driver/mongo"

type Model struct {
	processingExports *mongo.Collection
}
