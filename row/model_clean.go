package row

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m Model) Clean(
	ctx context.Context,
	eventID primitive.ObjectID,
) (
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	_, err = m.rows.DeleteMany(ctx, row{
		EventID: eventID,
	})
	return
}
