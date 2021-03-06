package resolvers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func shouldCount(ctx context.Context) bool {
	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		if field.Name == "total" {
			return true
		}
	}
	return false
}
