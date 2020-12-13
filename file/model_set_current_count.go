package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m Model) SetCurrentCount(ctx context.Context, eventID primitive.ObjectID, currentCount uint32) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = m.files.UpdateOne(ctx, file{
		EventID: eventID,
	}, bson.M{
		"$set": file{
			CurrentCount: currentCount,
		},
	})
	return
}
