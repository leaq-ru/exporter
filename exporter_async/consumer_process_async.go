package exporter_async

import (
	"encoding/json"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessAsync func(*parser.GetV2Request) error

func (c Consumer) ProcessAsync(req *parser.GetV2Request) (err error) {
	type msg struct {
		ID    primitive.ObjectID   `json:"i"`
		Query *parser.GetV2Request `json:"q"`
	}

	bytes, err := json.Marshal(msg{
		ID:    primitive.NewObjectID(),
		Query: req,
	})
	if err != nil {
		c.logger.Error().Err(err).Send()
		return
	}

	err = c.stanConn.Publish(exportSubjectName, bytes)
	if err != nil {
		c.logger.Error().Err(err).Send()
	}
	return
}
