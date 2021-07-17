package cached_export

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

var ErrNoFound = errors.New("cached url not found")

func (m Model) Get(ctx context.Context, query *parser.GetV2Request) (s3URL string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bytes, err := protojson.Marshal(query)
	if err != nil {
		return
	}

	sum := md5.Sum(bytes)

	var doc cachedExport
	err = m.cachedExports.FindOne(ctx, cachedExport{
		MD5: hex.EncodeToString(sum[:]),
	}).Decode(&doc)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		err = ErrNoFound
		return
	}

	s3URL = doc.URL
	return
}
