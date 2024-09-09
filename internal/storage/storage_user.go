package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) SaveUser(ctx context.Context, user *dto.UserEntity) (*dto.UserEntity, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Errorw("storage_user: error when begin *sqlx.TX when save user to DB", "error", err.Error())
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.NamedExecContext(
		ctx,
		"insert into users (login, password_hash, salt) values (:login, :password_hash, :salt)",
		user,
	)
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrEntityExists
		}

		log.Errorw("storage_user: unexpected error when save user to DB", "error", err.Error())
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorw("storage_user: error when commit *sqlx.TX when save user to DB", "error", err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) GetUserByLoginAndPasswordHash(ctx context.Context, login string, passwordHash string) (*dto.UserEntity, error) {
	sqlxRow := s.db.QueryRowxContext(
		ctx,
		"select * from users where login = $1 and password_hash = $2",
		login,
		passwordHash,
	)

	var userEntity dto.UserEntity
	if err := sqlxRow.StructScan(&userEntity); err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		log.Errorw(
			"storage_user: unexpected error when execute query",
			"error",
			err.Error(),
		)

		return nil, ErrUnexpextedDBError
	}

	return &userEntity, nil
}
