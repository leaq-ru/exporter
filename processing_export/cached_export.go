package processing_export

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type processingExport struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	EventID    primitive.ObjectID `bson:"e,omitempty"`
	LastActive time.Time          `bson:"l,omitempty"`
}
