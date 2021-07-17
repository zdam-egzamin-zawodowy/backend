package middleware

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/dataloader"
)

var (
	dataLoaderToContext contextKey = "data_loader"
)

func DataLoaderToContext(cfg dataloader.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dataLoaderToContext, dataloader.New(cfg))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DataLoaderFromContext(ctx context.Context) (*dataloader.DataLoader, error) {
	dataLoader := ctx.Value(dataLoaderToContext)
	if dataLoader == nil {
		err := errors.New("couldn't retrieve dataloader.DataLoader")
		return nil, err
	}

	dl, ok := dataLoader.(*dataloader.DataLoader)
	if !ok {
		err := errors.New("dataloader.DataLoader has wrong type")
		return nil, err
	}
	return dl, nil
}
