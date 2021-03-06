package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/zdam-egzamin-zawodowy/backend/internal/gin/middleware"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
	errorutils "github.com/zdam-egzamin-zawodowy/backend/pkg/utils/error"
)

type Directive struct{}

func (d *Directive) Authenticated(ctx context.Context, obj interface{}, next graphql.Resolver, yes bool) (interface{}, error) {
	_, err := middleware.UserFromContext(ctx)
	if yes && err != nil {
		return nil, errorutils.Wrap(err, messageMustBeSignedIn)
	} else if !yes && err == nil {
		return nil, fmt.Errorf(messageMustBeSignedOut)
	}

	return next(ctx)
}

func (d *Directive) HasRole(ctx context.Context, obj interface{}, next graphql.Resolver, role models.Role) (interface{}, error) {
	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		return nil, errorutils.Wrap(err, messageMustBeSignedIn)
	}
	if user.Role != role {
		return nil, fmt.Errorf(messageUnauthorized)
	}

	return next(ctx)
}
