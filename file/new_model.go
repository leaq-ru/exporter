package file

import (
	"github.com/nnqq/scr-exporter/mongo"
	m "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *m.Database) Model {
	return Model{
		db:    db,
		files: db.Collection(mongo.CollFile),
	}
}
