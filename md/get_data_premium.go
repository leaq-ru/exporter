package md

import (
	"context"
	"github.com/nnqq/scr-exporter/safeerr"
	"google.golang.org/grpc/metadata"
)

func GetDataPremium(ctx context.Context) (premium bool, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		err = safeerr.InternalServerError
		return
	}

	val := md.Get("data-premium")
	if len(val) != 0 {
		premium = val[0] == "true"
	}
	return
}
