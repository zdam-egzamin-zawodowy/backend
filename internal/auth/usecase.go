package auth

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/models"
)

type Usecase interface {
	SignIn(ctx context.Context, email, password string, staySignedIn bool) (*models.User, string, error)
	ExtractAccessTokenMetadata(ctx context.Context, accessToken string) (*models.User, error)
}
