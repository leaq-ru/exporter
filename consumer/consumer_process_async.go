package consumer

import (
	"encoding/json"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type message struct {
	ID    primitive.ObjectID     `json:"i"`
	Query *parser.GetListRequest `json:"q"`
}

type ProcessAsync func(primitive.ObjectID, *parser.GetListRequest) error

func (c Consumer) ProcessAsync(id primitive.ObjectID, req *parser.GetListRequest) (err error) {
	bytes, err := json.Marshal(message{
		ID:    id,
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
