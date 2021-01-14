package processing_export

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (m Model) SetProcessing(ctx context.Context, eventID primitive.ObjectID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = m.processingExports.UpdateOne(ctx, processingExport{
		EventID: eventID,
	}, bson.M{
		"$set": processingExport{
			LastActive: time.Now(),
		},
	}, options.Update().SetUpsert(true))
	return
}
