package file

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ErrDailyLimitExceeded = errors.New("daily limit exceeded. Try again tomorrow")

func (m Model) EnsureLimitAndCreatePending(ctx context.Context, userID, eventID primitive.ObjectID, name string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const dailyLimit = 30

	count, err := m.files.CountDocuments(ctx, file{
		UserID: userID,
	}, options.Count().SetLimit(dailyLimit))
	if err != nil {
		return
	}

	if count >= dailyLimit {
		return ErrDailyLimitExceeded
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
