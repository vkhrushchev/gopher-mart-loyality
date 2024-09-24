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
		log.Errorw("storage_user: error when begin *sqlx.TX", "error", err.Error())
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

		log.Errorw("storage_user: unexpected error when insert user to DB", "error", err.Error())
		return nil, err
	}

	userBalanceEntity := dto.UserBalanceEntity{
		UserLogin:          user.Login,
		TotalSum:           0.0,
		TotalWithdrawalSum: 0.0,
	}
	_, err = tx.NamedExecContext(
		ctx,
		"insert into user_balance (user_login, total_sum, total_withdrawal_sum) values (:user_login, :total_sum, :total_withdrawal_sum)",
		&userBalanceEntity,
	)
	if err != nil {
		log.Errorw("storage_user: unexpected storage error", "error", err.Error())
		return nil, ErrUnexpextedDBError
	}

	if err := tx.Commit(); err != nil {
		log.Errorw("storage_user: error when commit *sqlx.TX", "error", err.Error())
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
	if err := sqlxRow.StructScan(&userEntity); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		log.Errorw(
			"storage_user: unexpected error when execute query",
			"error",
			err.Error())

		return nil, ErrUnexpextedDBError
	}

	return &userEntity, nil
}

func (s *UserStorage) GetUserBalanceByLogin(ctx context.Context, login string) (*dto.UserBalanceEntity, error) {
	sqlxRow := s.db.QueryRowxContext(
		ctx,
		"select * from user_balance where user_login = $1",
		login,
	)

	var userBalanceEntity dto.UserBalanceEntity
	if err := sqlxRow.StructScan(&userBalanceEntity); errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		log.Errorw(
			"storage_user: unexpected error when execute query",
			"error",
			err.Error())

		return nil, ErrUnexpextedDBError
	}

	return &userBalanceEntity, nil
}
