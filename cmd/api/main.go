package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/askme/api/internal/config"
	"github.com/askme/api/pkg/arango"
	"github.com/askme/api/pkg/middleware"
	"github.com/askme/api/pkg/slogutil"
)

func main() {
	logger := slog.New(slogutil.NewIndentedJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize ArangoDB client
	db, err := arango.NewClient(cfg.ArangoDB)
	if err != nil {
		slog.Error("failed to connect to ArangoDB", "error", err)
		os.Exit(1)
	}

	// Initialize app with dependency injection
	app := NewApp(db)

	// Setup router
	mux := http.NewServeMux()
	app.RegisterRoutes(mux)

	// Apply middleware chain (order matters: outermost first)
	// Recovery -> RequestID -> Logger -> SecureHeaders -> CORS -> FakeAuth -> JSON -> handler
	handler := middleware.Chain(
		mux,
		middleware.Recovery,      // Recover from panics (outermost)
		middleware.RequestID,     // Add request ID for tracing
		middleware.Logger,        // Log all requests
		middleware.SecureHeaders, // Add security headers
		middleware.CORS(middleware.DefaultCORSConfig()),         // Handle CORS
		middleware.FakeAuth(middleware.DefaultFakeAuthConfig()), // Fake auth for dev (extracts X-User-ID header)
		middleware.JSON, // Set JSON content type
	)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}
