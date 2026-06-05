package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ready4god2513/requestbin/handler"
	"github.com/ready4god2513/requestbin/store"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://requestbin:requestbin@localhost:5432/requestbin?sslmode=disable"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}
	defer pool.Close()

	if err := store.Migrate(ctx, pool); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	st := store.New(pool)
	hub := handler.NewHub()
	go hub.Run()

	h := handler.New(st, hub)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization, "X-Requested-With",
		},
		AllowMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut,
			http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodHead,
		},
	}))

	api := e.Group("/api")
	api.POST("/bins", h.CreateBin)
	api.GET("/bins/:id", h.GetBin)
	api.DELETE("/bins/:id", h.DeleteBin)
	api.GET("/bins/:id/requests", h.ListRequests)
	api.DELETE("/bins/:id/requests", h.ClearRequests)
	api.DELETE("/bins/:id/requests/:reqID", h.DeleteRequest)
	api.GET("/bins/:id/sse", h.SSE)

	e.Any("/r/:id", h.Capture)
	e.Any("/r/:id/*", h.Capture)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-sctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
