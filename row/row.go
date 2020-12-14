package row

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type row struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	EventID       primitive.ObjectID `bson:"e,omitempty"`
	FullCompanyV2 []byte             `bson:"f,omitempty"`
	CreatedAt     time.Time          `bson:"ca,omitempty"`
}
