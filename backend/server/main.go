package main

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"support/internal/config"
	"support/internal/db/postgres"
	"support/internal/handlers"
	"support/internal/repository"
	"support/internal/usecase"
	"syscall"
	"time"
)

func main() {
	err := config.SetUp()
	if err != nil {
		slog.Error("failed to fetch config", "error", err)
	}
	postgresDB, err := postgres.InitDB()
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
	}
	err = postgres.MakeMigrations(true)
	if err != nil {
		slog.Error("failed to make migrations", "error", err)
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	knowledgeRepo := repository.NewKnowledgeRepo(postgresDB)
	KnowledgeUseCase := usecase.NewKnowledgeInstance(knowledgeRepo)
	knowledgeHandler := handlers.NewKnowledgeInstance(KnowledgeUseCase)

	e.GET("/api/answer", knowledgeHandler.Answer)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := e.Start(":" + config.AppConfig.Server.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	<-stop
	slog.Info("received shutdown signal, starting shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		slog.Error("failed to gracefully shut down server", "error", err)
	}

	slog.Info("server gracefully stopped")

}
