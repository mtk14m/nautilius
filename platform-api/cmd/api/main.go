package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mtk14m/notilius/platform-api/internal/config"
	"github.com/mtk14m/notilius/platform-api/internal/middleware"
)

func main() {

	//Load config
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	//Setup du looger avec la config
	setupLogger(cfg.Observability.LogLevel, cfg.Observability.LogFormat)
	slog.Info("Starting platform-api", "version", "1.0.0")

	//Router
	router := gin.Default()

	// Middleware: Error handler (must be first to catch all errors)
	router.Use(middleware.ErrorHandlerMiddleware())

	// Middleware: Request logging
	router.Use(requestLoggingMiddleware())

	//Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "platform-api",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "1.0.0",
				"service": "platform-api",
			})
		})
	}

	//SERVER
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	//Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Server started", "addr", srv.Addr)

	//Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}

// setupLogger configures slog based on config
func setupLogger(level, format string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: logLevel,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	slog.SetDefault(slog.New(handler))
}

// requestLoggingMiddleware logs HTTP requests
func requestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate trace ID if not present
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = randomString(16)
		}
		c.Set("traceId", traceID)
		c.Header("X-Trace-ID", traceID)

		// Process request
		c.Next()

		// Log request
		duration := time.Since(start)
		slog.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"trace_id", traceID,
			"ip", c.ClientIP(),
		)
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
