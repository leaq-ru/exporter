package row

import (
	"context"
	"time"
)

func (m Model) Flush(
	ctx context.Context,
) (
	err error,
) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(m.state.buf) != 0 {
		_, err = m.rows.InsertMany(ctx, m.state.buf)
		if err != nil {
			return
		}

		m.state.buf = []interface{}{}
	}
	return
}
