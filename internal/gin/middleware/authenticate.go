package middleware

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
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

func UserFromContext(ctx context.Context) (*models.User, error) {
	user := ctx.Value(authenticateKey)
	if user == nil {
		err := fmt.Errorf("Could not retrieve *models.User")
		return nil, err
	}

	u, ok := user.(*models.User)
	if !ok {
		err := fmt.Errorf("*models.User has wrong type")
		return nil, err
	}
	return u, nil
}
