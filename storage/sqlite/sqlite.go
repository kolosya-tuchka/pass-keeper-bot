package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"pass-keeper-bot/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, userService *storage.UserService) error {
	q := `INSERT INTO services (userID, service, login, password) VALUES ($1, $2, $3, $4)
 		  ON CONFLICT(userID, service) DO UPDATE SET login=$3, password=$4`

	if _, err := s.db.ExecContext(
		ctx, q,
		userService.UserID,
		userService.Service,
		userService.Login,
		userService.Password,
	); err != nil {
		return fmt.Errorf("can't save service: %w", err)
	}

	return nil
}

func (s *Storage) Pick(ctx context.Context, userID int64, service string) (*storage.UserService, error) {
	q := `SELECT login, password FROM services WHERE userID=? AND service=?`

	var login, password string

	err := s.db.QueryRowContext(
		ctx, q,
		userID,
		service,
	).Scan(&login, &password)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSuchService
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick service: %w", err)
	}

	return &storage.UserService{UserID: userID, Service: service, Login: login, Password: password}, nil
}

func (s *Storage) Remove(ctx context.Context, userID int64, service string) error {
	q := `DELETE FROM services WHERE userID=? AND service=?`
	res, err := s.db.ExecContext(
		ctx, q,
		userID,
		service,
	)

	if err != nil {
		return fmt.Errorf("can't remove service: %w", err)
	}

	count, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("can't take result: %w", err)
	}

	if count == 0 {
		return storage.ErrNoSuchService
	}

	return nil
}

func (s *Storage) Exists(ctx context.Context, userID int64, service string) (bool, error) {
	q := `SELECT COUNT(*) FROM services WHERE userID=? AND service=?`

	var count int

	if err := s.db.QueryRowContext(ctx, q,
		userID,
		service).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if service exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS
    services (userID BIGINT, service TEXT, login TEXT, password TEXT, PRIMARY KEY(userID, service))`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
