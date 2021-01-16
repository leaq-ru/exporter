package csv

import (
	"context"
	"errors"
	"github.com/nnqq/scr-proto/codegen/go/parser"
	"golang.org/x/sync/errgroup"
	"io"
)

func DoPipelineSync(
	ctx context.Context,
	compStream parser.Company_GetFullClient,
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
			}
		}
	})

	eg.Go(func() (e error) {
		defer csvCtxCancel()
		csvPath, e = Create(compCh)
		return
	})
	err = eg.Wait()
	return
}
