package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"

	pm "github.com/stanistan/present-me"
	"github.com/stanistan/present-me/internal/cache"
)

func main() {
	var config pm.Config
	_ = kong.Parse(&config)

	config.Configure()
	log := config.Logger()

	gh, err := config.GH()
	if err != nil {
		log.Fatal().Err(err).Msg("could not configure GH client")
	}

	r := mux.NewRouter()

	// 1. Register API routes & middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(
		hlog.NewHandler(log),
		githubMiddleware(gh),
		cacheMiddleware,
	)

	for _, r := range apiRoutes {
		api.
			Handle(r.Prefix, r.Handler).
			Methods(r.Method)
	}

	// 2. Register fallback website handler
	websiteHandler, err := config.WebsiteHandler()
	if err != nil {
		log.Fatal().Err(err).Msg("could not build handler")
	}
	r.PathPrefix("/").Handler(websiteHandler)

	// 3. Init server
	s := &http.Server{
		Addr:         config.Address(),
		ReadTimeout:  config.ServerReadTimeout,
		WriteTimeout: config.ServerWriteTimeout,
		Handler:      r,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

type ghCtxKey int

var ghCtxValue ghCtxKey

func ContextWithGH(ctx context.Context, gh *pm.GH) context.Context {
	return context.WithValue(ctx, ghCtxValue, gh)
}

func GHFromContext(ctx context.Context) (*pm.GH, bool) {
	v, ok := ctx.Value(ghCtxValue).(*pm.GH)
	return v, ok
}

func githubMiddleware(g *pm.GH) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithGH(r.Context(), g)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := cache.ContextWithOptions(r.Context(), &cache.Options{
			TTL:          10 * time.Minute,
			ForceRefresh: r.URL.Query().Get("refresh") == "1",
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
