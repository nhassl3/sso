package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
	"github.com/nhassl3/sso-app/internals/domain/models"
	"github.com/nhassl3/sso-app/internals/lib/logger/sl"
	"github.com/nhassl3/sso-app/internals/storage"
)

const (
	opNewStorage = "storage.sqlite.NewStorage"
	opSaveUser   = "storage.sqlite.SaveUser"
	opUser       = "storage.sqlite.User"
	opIsAdmin    = "storage.sqlite.IsAdmin"
	opApp        = "storage.sqlite.App"
)

type Storage struct {
	db *sql.DB
}

// NewStorage creates a new instance of the SQLite storage.
func NewStorage(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, sl.ErrUpLevel(opNewStorage, err.Error())
	}

	return &Storage{db: db}, nil
}

// SaveUser save user in the system
func (s *Storage) SaveUser(ctx context.Context, email string, hashPassword []byte) (userID int64, err error) {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, sl.ErrUpLevel(opSaveUser, err.Error())
	}
	defer stmt.Close()

	var sqliteErr sqlite3.Error

	res, err := stmt.ExecContext(ctx, email, hashPassword)
	if err != nil {
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, sl.ErrUpLevel(opSaveUser, storage.ErrUserExists.Error())
		}

		return 0, sl.ErrUpLevel(opSaveUser, err.Error())
	}

	userID, err = res.LastInsertId()
	if err != nil {
		return 0, sl.ErrUpLevel(opSaveUser, err.Error())
	}

	return
}

// User returns user model by email
func (s *Storage) User(ctx context.Context, email string) (user models.User, err error) {
	err = s.newSelect(
		ctx,
		"SELECT id, email, pass_hash FROM users WHERE email=?",
		[]interface{}{email},
		&user.ID, &user.Email, &user.HashPassword,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, sl.ErrUpLevel(opUser, storage.ErrUserNotFound.Error())
		}

		return models.User{}, sl.ErrUpLevel(opUser, err.Error())
	}

	return
}

// IsAdmin checks by UID if user is admin returns true else false
func (s *Storage) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	err = s.newSelect(
		ctx,
		`SELECT EXISTS(
    SELECT 1 FROM admins
             WHERE admins.user_id = ?
             AND EXISTS(SELECT 1 FROM users WHERE users.id = ?)
)`,
		[]interface{}{userID, userID}, // Two user IDs need to be transferred
		&isAdmin,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opIsAdmin, storage.ErrUserNotFound.Error())
		}

		return false, sl.ErrUpLevel(opIsAdmin, err.Error())
	}

	return
}

// App returns model of an application
func (s *Storage) App(ctx context.Context, appID int32) (app models.App, err error) {
	err = s.newSelect(
		ctx,
		"SELECT id, name, secret FROM apps WHERE id = ?",
		[]interface{}{appID},
		&app.ID, &app.Name, &app.Secret,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, sl.ErrUpLevel(opApp, storage.ErrAppNotFound.Error())
		}

		return models.App{}, sl.ErrUpLevel(opApp, err.Error())
	}

	return
}

// newSelect cleaning code deletes duplicates
func (s *Storage) newSelect(ctx context.Context, query string, args []interface{}, dest ...interface{}) error {
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)

	err = row.Scan(dest...)
	if err != nil {
		return err
	}

	return nil
}
