package httpdelivery

import (
	"fmt"
	"github.com/Kichiyaki/appmode"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/querycomplexity"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/directive"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/resolvers"
)

const (
	playgroundTTL      = time.Hour / time.Second
	graphqlEndpoint    = "/graphql"
	playgroundEndpoint = "/"
)

type Config struct {
	Resolver  *resolvers.Resolver
	Directive *directive.Directive
}

func Attach(r chi.Router, cfg Config) error {
	if cfg.Resolver == nil {
		return errors.New("cfg.Resolver is required")
	}
	gqlHandler := graphqlHandler(prepareConfig(cfg.Resolver, cfg.Directive))
	r.Get(graphqlEndpoint, gqlHandler)
	r.Post(graphqlEndpoint, gqlHandler)
	if appmode.Equals(appmode.DevelopmentMode) {
		r.Get(playgroundEndpoint, playgroundHandler())
	}
	return nil
}

func graphqlHandler(cfg generated.Config) http.HandlerFunc {
	srv := handler.New(generated.NewExecutableSchema(cfg))

	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: 32 << 18,
		MaxMemory:     32 << 18,
	})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.SetQueryCache(lru.New(100))
	srv.Use(querycomplexity.GetComplexityLimitExtension())
	if appmode.Equals(appmode.DevelopmentMode) {
		srv.Use(extension.Introspection{})
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Cache-Control", "no-store, must-revalidate")
		srv.ServeHTTP(w, r)
	}
}

func playgroundHandler() http.HandlerFunc {
	h := playground.Handler("Playground", graphqlEndpoint)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Cache-Control", fmt.Sprintf(`public, max-age=%d`, playgroundTTL))
		h.ServeHTTP(w, r)
	}
}

func prepareConfig(r *resolvers.Resolver, d *directive.Directive) generated.Config {
	cfg := generated.Config{
		Resolvers:  r,
		Complexity: querycomplexity.GetComplexityRoot(),
	}
	cfg.Directives.Authenticated = d.Authenticated
	cfg.Directives.HasRole = d.HasRole
	return cfg
}
