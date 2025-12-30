package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	GRPCS *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
	db *pgxpool.Pool,
) *App {
	storage := postgres.New(db)

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCS: grpcApp,
	}
}
