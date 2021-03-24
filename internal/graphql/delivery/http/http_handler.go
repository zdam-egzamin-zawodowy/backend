package httpdelivery

import (
	"fmt"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/querycomplexity"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/directive"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/generated"
	"github.com/zdam-egzamin-zawodowy/backend/internal/graphql/resolvers"
	"github.com/zdam-egzamin-zawodowy/backend/pkg/mode"
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

func Attach(group *gin.RouterGroup, cfg Config) error {
	if cfg.Resolver == nil {
		return fmt.Errorf("Graphql resolver cannot be nil")
	}
	gqlHandler := graphqlHandler(prepareConfig(cfg.Resolver, cfg.Directive))
	group.GET(graphqlEndpoint, gqlHandler)
	group.POST(graphqlEndpoint, gqlHandler)
	if mode.Get() == mode.DevelopmentMode {
		group.GET(playgroundEndpoint, playgroundHandler())
	}
	return nil
}

// Defining the GraphQL handler
func graphqlHandler(cfg generated.Config) gin.HandlerFunc {
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
	if mode.Get() == mode.DevelopmentMode {
		srv.Use(extension.Introspection{})
	}

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, must-revalidate")
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("Playground", graphqlEndpoint)

	return func(c *gin.Context) {
		c.Header("Cache-Control", fmt.Sprintf(`public, max-age=%d`, playgroundTTL))
		h.ServeHTTP(c.Writer, c.Request)
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
