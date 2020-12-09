package file

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type file struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"u,omitempty"`
	EventID      primitive.ObjectID `bson:"e,omitempty"`
	Name         string             `bson:"n,omitempty"`
	URL          string             `bson:"ur,omitempty"`
	Status       status             `bson:"s,omitempty"`
	CurrentCount uint32             `bson:"c,omitempty"`
	TotalCount   uint32             `bson:"t,omitempty"`
	CreatedAt    time.Time          `bson:"ca,omitempty"`
}

type status uint8

const (
	_ status = iota
	status_pending
	status_inProgress
	status_success
	status_fail
)
