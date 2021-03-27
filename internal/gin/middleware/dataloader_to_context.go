package middleware

import (
	"context"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/dataloader"
)

var (
	dataLoaderToContext contextKey = "data_loader"
)

func DataLoaderToContext(cfg dataloader.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), dataLoaderToContext, dataloader.New(cfg))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func DataLoaderFromContext(ctx context.Context) (*dataloader.DataLoader, error) {
	dataLoader := ctx.Value(dataLoaderToContext)
	if dataLoader == nil {
		err := errors.New("could not retrieve dataloader.DataLoader")
		return nil, err
	}

	dl, ok := dataLoader.(*dataloader.DataLoader)
	if !ok {
		err := errors.New("dataloader.DataLoader has wrong type")
		return nil, err
	}
	return dl, nil
}
