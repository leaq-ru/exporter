package file

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m Model) WatchJob(
	ctx context.Context,
	masterEventID primitive.ObjectID,
	slaveEventID primitive.ObjectID,
) (
	err error,
) {
	for {
		var master file
		err = m.files.FindOne(ctx, file{
			EventID: masterEventID,
		}).Decode(&master)
		if err != nil {
			return
		}

		switch master.Status {
		case status_inProgress:
			_, err = m.files.UpdateOne(ctx, file{
				EventID: slaveEventID,
			}, bson.M{
				"$set": file{
					CurrentCount: master.CurrentCount,
					TotalCount:   master.TotalCount,
				},
			})
			if err != nil {
				return
			}
		case status_success:
			_, err = m.files.UpdateOne(ctx, file{
				EventID: slaveEventID,
			}, bson.M{
				"$set": file{
					URL: master.URL,
				},
			})
			return
		case status_fail:
			_, err = m.files.UpdateOne(ctx, file{
				EventID: slaveEventID,
			}, bson.M{
				"$set": file{
					Status: status_fail,
				},
			})
			if err != nil {
				return
			}

			err = errors.New("master job failed")
			return
		}

		time.Sleep(10 * time.Second)
	}
}
