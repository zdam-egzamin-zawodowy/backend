package auth

import (
	"context"
	"github.com/zdam-egzamin-zawodowy/backend/internal"
)

type Usecase interface {
	SignIn(ctx context.Context, email, password string, staySignedIn bool) (*internal.User, string, error)
	ExtractAccessTokenMetadata(ctx context.Context, accessToken string) (*internal.User, error)
}
