package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

type WithdrawlStorage struct {
	db *sqlx.DB
}

func NewWithdrawalStorage(db *sqlx.DB) *WithdrawlStorage {
	return &WithdrawlStorage{
		db: db,
	}
}

func (s *WithdrawlStorage) SaveBalanceWithdrawal(ctx context.Context, balanceWithdraw *dto.BalanceWithdrawalEntity) (*dto.BalanceWithdrawalEntity, error) {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Errorw("storage_withdrawl: error when begin *sqlx.TX", "error", err.Error())
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	sqlxRow := tx.QueryRowxContext(ctx, "select * from user_balance where user_login = $1 for update", balanceWithdraw.UserLogin)
	if sqlxRow.Err() != nil {
		log.Errorw("storage_withdrawl: unexpected db error", "error", sqlxRow.Err())
		return nil, ErrUnexpextedDBError
	}

	var userBalance dto.UserBalanceEntity
	if err := sqlxRow.StructScan(&userBalance); err != nil {
		log.Errorw("storage_withdrawl: unexpected db error", "error", err.Error())
		return nil, ErrUnexpextedDBError
	}

	// повторно проверяем баланс после блокировки строки в БД
	if userBalance.TotalSum < balanceWithdraw.WithdrawalSum {
		return nil, ErrNoFundsOnBalance
	}

	userBalance.TotalSum = userBalance.TotalSum - balanceWithdraw.WithdrawalSum
	userBalance.TotalWithdrawalSum = userBalance.TotalWithdrawalSum + balanceWithdraw.WithdrawalSum

	_, err = tx.ExecContext(
		ctx,
		"update user_balance set total_sum = $1, total_withdrawal_sum = $2 where id = $3",
		userBalance.TotalSum, userBalance.TotalWithdrawalSum, userBalance.ID)
	if err != nil {
		log.Errorw("storage_withdrawl: unexpected db error", "error", err.Error())
		return nil, ErrUnexpextedDBError
	}

	_, err = tx.NamedExecContext(
		ctx,
		"insert into balance_withdrawals (user_login, order_number, withdrawal, processed_at) values (:user_login, :order_number, :withdrawal, :processed_at)",
		balanceWithdraw)
	if err != nil {
		log.Errorw("storage_withdrawl: unexpected db error", "error", err.Error())
		return nil, ErrUnexpextedDBError
	}

	if err := tx.Commit(); err != nil {
		log.Errorw("storage_withdrawl: error when commit *sqlx.TX", "error", err.Error())
		return nil, err
	}

	return balanceWithdraw, nil
}

func (s *WithdrawlStorage) GetBalanceWithdrawalsByUserLogin(ctx context.Context, userLogin string) ([]dto.BalanceWithdrawalEntity, error) {
	sqlxRows, err := s.db.QueryxContext(ctx, "select * from balance_withdrawals where user_login = $1", userLogin)
	if err != nil {
		log.Errorw("storage_withdrawl: unexpected db error", "error", err.Error())
		return nil, ErrUnexpextedDBError
	}

	balanceWithdrawalEntities := make([]dto.BalanceWithdrawalEntity, 0)
	for sqlxRows.Next() {
		var balanceWithdrawalEntity dto.BalanceWithdrawalEntity
		if err := sqlxRows.StructScan(&balanceWithdrawalEntity); err != nil {
			log.Errorw("storage_withdrawl: unexpected db error", "error", err.Error())
			return nil, ErrUnexpextedDBError
		}

		balanceWithdrawalEntities = append(balanceWithdrawalEntities, balanceWithdrawalEntity)
	}

	return balanceWithdrawalEntities, nil
}
