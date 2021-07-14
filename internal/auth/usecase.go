package auth

import (
	"context"

	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

type Usecase interface {
	SignIn(ctx context.Context, email, password string, staySignedIn bool) (*model.User, string, error)
	ExtractAccessTokenMetadata(ctx context.Context, accessToken string) (*model.User, error)
}
