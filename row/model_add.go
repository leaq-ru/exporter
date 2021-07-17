package row

import (
	"context"
	"github.com/leaq-ru/proto/codegen/go/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"time"
)

func (m Model) Add(
	ctx context.Context,
	eventID primitive.ObjectID,
	fullCompanyV2 *parser.FullCompanyV2,
) (
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	bytes, err := proto.Marshal(fullCompanyV2)
	if err != nil {
		return
	}

	m.state.mu.Lock()
	m.state.buf = append(m.state.buf, row{
		EventID:       eventID,
		FullCompanyV2: bytes,
		CreatedAt:     time.Now().UTC(),
	})
	m.state.mu.Unlock()

	if len(m.state.buf) >= 300 {
		_, err = m.rows.InsertMany(ctx, m.state.buf)
		if err != nil {
			return
		}

		m.state.buf = []interface{}{}
	}
	return
}
