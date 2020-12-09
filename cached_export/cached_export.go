package cached_export

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type cachedExport struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	MD5       string             `bson:"m,omitempty"`
	URL       string             `bson:"u,omitempty"`
	CreatedAt time.Time          `bson:"ca,omitempty"`
}
