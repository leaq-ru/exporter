package file

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

func (m Model) GetMasterJob(
	ctx context.Context,
	query *parser.GetV2Request,
	selfEventID primitive.ObjectID,
) (
	eventID primitive.ObjectID,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	md5key, err := makeMD5Key(query)
	if err != nil {
		return
	}

	var doc file
	err = m.files.FindOne(ctx, bson.M{
		"_id": bson.M{
			"$ne": selfEventID,
		},
		"m": md5key,
		"s": status_inProgress,
	}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = nil
		}
		return
	}

	eventID = doc.EventID
	return
}

func makeMD5Key(query *parser.GetV2Request) (key string, err error) {
	bytes, err := protojson.Marshal(query)
	if err != nil {
		return
	}

	sum := md5.Sum(bytes)
	key = hex.EncodeToString(sum[:])
	return
}
