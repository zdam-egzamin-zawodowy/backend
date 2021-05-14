package middleware

import (
	"context"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

var (
	ginContextToContextKey contextKey = "gin_context"
)

func GinContextToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextToContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextToContextKey)
	if ginContext == nil {
		err := errors.New("couldn't retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := errors.New("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
