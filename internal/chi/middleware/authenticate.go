package middleware

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	"github.com/zdam-egzamin-zawodowy/backend/internal/auth"
	"github.com/zdam-egzamin-zawodowy/backend/internal/model"
)

const (
	authorizationHeader = "Authorization"
)

var (
	authenticateKey contextKey = "current_user"
)

func Authenticate(ucase auth.Usecase) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r.Header.Get("Authorization"))
			if token != "" {
				ctx := r.Context()
				user, err := ucase.ExtractAccessTokenMetadata(ctx, token)
				if err == nil && user != nil {
					ctx = context.WithValue(ctx, authenticateKey, user)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
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
