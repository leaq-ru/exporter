package processing_export

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m Model) UnsetProcessing(ctx context.Context, eventID primitive.ObjectID) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = m.processingExports.DeleteOne(ctx, processingExport{
		EventID: eventID,
	})
	return
}
