package file

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ErrConcExports = errors.New("too many concurrent exports. Wait for old export succeeded, and try again")

func (m Model) EnsureLimitAndCreatePending(ctx context.Context, userID, eventID primitive.ObjectID, name string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const concLimit = 3

	count, err := m.files.CountDocuments(ctx, bson.M{
		"u": userID,
		"s": bson.M{
			"$in": bson.A{
				status_pending,
				status_inProgress,
			},
		},
	}, options.Count().SetLimit(concLimit))
	if err != nil {
		return
	}

	if count >= concLimit {
		return ErrConcExports
	}

	_, err = m.files.InsertOne(ctx, file{
		UserID:    userID,
		EventID:   eventID,
		Name:      name,
		Status:    status_pending,
		CreatedAt: time.Now().UTC(),
	})
	return
}
