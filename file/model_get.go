package file

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (m Model) Get(
	ctx context.Context,
	userID primitive.ObjectID,
	skip,
	limit uint32,
) (
	res []file,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cur, err := m.files.Find(ctx, file{
		UserID: userID,
	}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return
	}

	err = cur.All(ctx, &res)
	return
}
