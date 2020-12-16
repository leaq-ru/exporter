package file

import (
	"context"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m Model) SetMasterJobInProgress(
	ctx context.Context,
	eventID primitive.ObjectID,
	query *parser.GetV2Request,
	totalCount uint32,
) (
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	md5key, err := makeMD5Key(query)
	if err != nil {
		return
	}

	_, err = m.files.UpdateOne(ctx, file{
		EventID: eventID,
	}, bson.M{
		"$set": file{
			MD5:        md5key,
			Status:     status_inProgress,
			TotalCount: totalCount,
		},
	})
	return
}
