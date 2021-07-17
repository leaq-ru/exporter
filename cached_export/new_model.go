package cached_export

import (
	"github.com/leaq-ru/exporter/mongo"
	m "go.mongodb.org/mongo-driver/mongo"
)

func NewModel(db *m.Database) Model {
	return Model{
		db:            db,
		cachedExports: db.Collection(mongo.CollCachedExport),
	}
}
