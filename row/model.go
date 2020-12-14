package row

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type state struct {
	mu  *sync.Mutex
	buf []interface{} // []row
}

type Model struct {
	db    *mongo.Database
	rows  *mongo.Collection
	state *state
}
