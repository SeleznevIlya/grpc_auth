package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SeleznevIlya/grpc_auth/internal/domain/models"
	"github.com/SeleznevIlya/grpc_auth/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of the PostgreSQL storage.
func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.postgres.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, password) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var postgresErr *pq.Error

		if errors.As(err, &postgresErr) && postgresErr.Code == pq.ErrorCode("23505") {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, email)

	var user models.User
	if err = res.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM is_admin WHERE user_id=?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	if err = res.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM app WHERE id=?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, appID)

	var app models.App
	if err = res.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
