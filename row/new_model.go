package row

import (
	"github.com/nnqq/scr-exporter/mongo"
	m "go.mongodb.org/mongo-driver/mongo"
	"sync"
)

func NewModel(db *m.Database) Model {
	return Model{
		db:   db,
		rows: db.Collection(mongo.CollRow),
		state: &state{
			mu:  &sync.Mutex{},
			buf: []interface{}{},
		},
	}
}
