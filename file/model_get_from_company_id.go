package file

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (m Model) GetFromCompanyID(ctx context.Context, eventID primitive.ObjectID) (companyID string, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var doc file
	err = m.files.FindOne(ctx, file{
		EventID: eventID,
	}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = nil
	}

	companyID = doc.FromCompanyID
	return
}
