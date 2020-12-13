package csv

import (
	"context"
	"errors"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"golang.org/x/sync/errgroup"
	"io"
	"sync/atomic"
)

func DoPipeline(
	ctx context.Context,
	compStream parser.Company_GetFullClient,
	currentCount *uint32,
) (
	csvPath string,
	err error,
) {
	csvCtx, csvCtxCancel := context.WithCancel(ctx)
	defer csvCtxCancel()

	var eg errgroup.Group
	compCh := make(chan *parser.FullCompanyV2)
	eg.Go(func() (e error) {
		defer close(compCh)

		for {
			comp, er := compStream.Recv()
			e = er
			if e != nil {
				if errors.Is(e, io.EOF) {
					e = nil
				}
				return
			}

			select {
			case <-csvCtx.Done():
				return
			default:
				compCh <- comp
				atomic.AddUint32(currentCount, 1)
			}
		}
	})

	eg.Go(func() (e error) {
		defer csvCtxCancel()
		csvPath, e = create(compCh)
		return
	})
	err = eg.Wait()
	return
}
