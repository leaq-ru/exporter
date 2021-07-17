package file

import (
	"github.com/leaq-ru/exporter/mongo"
	m "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *m.Database) Model {
	return Model{
		db:    db,
		files: db.Collection(mongo.CollFile),
	}
}
