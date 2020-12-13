package cached_export

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

func (m Model) Set(ctx context.Context, query *parser.GetV2Request, s3URL string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bytes, err := protojson.Marshal(query)
	if err != nil {
		return
	}

	sum := md5.Sum(bytes)

	_, err = m.cachedExports.UpdateOne(ctx, cachedExport{
		MD5: hex.EncodeToString(sum[:]),
	}, bson.M{
		"$set": cachedExport{
			URL:       s3URL,
			CreatedAt: time.Now().UTC(),
		},
	}, options.Update().SetUpsert(true))
	return
}
