package processing_export

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (m Model) IsProcessing(ctx context.Context, eventID primitive.ObjectID) (processing bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = m.processingExports.FindOne(ctx, processingExport{
		EventID: eventID,
	}).Err()
	if err == nil {
		processing = true
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = nil
		return
	}
	return
}
