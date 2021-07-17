package row

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/leaq-ru/exporter/csv"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"time"
)

func (m Model) PipeMongoToLocalFile(
	ctx context.Context,
	eventID primitive.ObjectID,
) (
	csvPath string,
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	cur, err := m.rows.Find(ctx, row{
		EventID: eventID,
	})
	if err != nil {
		return
	}

	ch := make(chan *parser.FullCompanyV2)
	var eg errgroup.Group
	eg.Go(func() (e error) {
		defer close(ch)

		for cur.Next(ctx) {
			var doc row
			e = cur.Decode(&doc)
			if e != nil {
				return
			}

			comp := &parser.FullCompanyV2{}
			e = proto.Unmarshal(doc.FullCompanyV2, comp)
			if e != nil {
				return
			}

			ch <- comp
		}

		e = cur.Err()
		return
	})

	eg.Go(func() (e error) {
		csvPath, e = csv.Create(ch)
		return
	})
	err = eg.Wait()
	return
}
