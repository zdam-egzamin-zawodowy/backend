package middleware

import (
	"context"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

const (
	authorizationHeader = "Authorization"
)

var (
	authenticateKey contextKey = "current_user"
)

func Authenticate(ucase auth.Usecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.GetHeader(authorizationHeader))
		if token != "" {
			ctx := c.Request.Context()
			user, err := ucase.ExtractAccessTokenMetadata(ctx, token)
			if err == nil && user != nil {
				ctx = context.WithValue(ctx, authenticateKey, user)
				c.Request = c.Request.WithContext(ctx)
			}
		}
		c.Next()
	}
}

func UserFromContext(ctx context.Context) (*model.User, error) {
	user := ctx.Value(authenticateKey)
	if user == nil {
		err := errors.New("couldn't retrieve *model.User")
		return nil, err
	}

	u, ok := user.(*model.User)
	if !ok {
		err := errors.New("*model.User has wrong type")
		return nil, err
	}
	return u, nil
}
