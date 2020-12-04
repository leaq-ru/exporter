package event_log

import (
	"github.com/nnqq/scr-exporter/mongo"
	m "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *m.Database) Model {
	return Model{
		db:        db,
		eventLogs: db.Collection(mongo.CollEventLog),
	}
}
