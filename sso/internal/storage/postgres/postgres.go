package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "postgres.SaveUser"
	var id int64

	if err := s.db.QueryRow(
		ctx,
		"INSERT INTO users(email,pass_hash) VALUES($1,$2) RETURNING id", email, passHash).
		Scan(&id); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { //Тобто такий запис існує
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
	}
	return id, nil
}
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.User"
	var user models.User

	if err := s.db.QueryRow(ctx, "SELECT id,email,pass_hash FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "postgres.IsAmin"
	var isAdmin bool

	if err := s.db.QueryRow(ctx, "SELECT isAdmin WHERE id=$1", userId).
		Scan(&isAdmin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "postgres.App"
	var app models.App

	if err := s.db.QueryRow(
		ctx,
		"SELECT name id secret FROM apps WHERE id=$1", appID,
	).Scan(&app.Name, &app.ID, &app.Secret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
