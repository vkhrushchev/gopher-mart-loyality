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

type OrderStorage struct {
	db *sqlx.DB
}

func NewOrderStorage(db *sqlx.DB) *OrderStorage {
	return &OrderStorage{
		db: db,
	}
}

func (s *OrderStorage) SaveOrder(ctx context.Context, order *dto.OrderEntity) (*dto.OrderEntity, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Errorw("storage_order: error when begin *sqlx.TX when save order to DB", "error", err.Error())
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.NamedExecContext(
		ctx,
		`
		insert into orders (user_login, number, accrual, status, uploaded_at) 
		values (:user_login, :number, :accrual, :status, :uploaded_at);
		`,
		order,
	)
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrEntityExists
		}

		log.Errorw("storage_order: unexpected error when save order to DB", "error", err.Error())
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorw("storage_order: error when commit *sqlx.TX when save order to DB", "error", err.Error())
		return nil, err
	}

	return order, nil
}

func (s *OrderStorage) UpdateOrderStatus(ctx context.Context, orderNumber string, orderStatus dto.OrderStatus) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Errorw("storage_order: error when begin tx", "error", err.Error())
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, "update orders set status = $1 where number = $2", orderStatus, orderNumber)
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	return nil
}

func (s *OrderStorage) UpdateOrderStatusAndAccrual(ctx context.Context, orderNumber string, orderStatus dto.OrderStatus, accrual float64) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Errorw("storage_order: error when begin tx", "error", err.Error())
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	sqlxRow := tx.QueryRowxContext(ctx, "select * from orders where number = $1", orderNumber)
	if sqlxRow.Err() != nil {
		log.Errorw(
			"storage_order: unexpected db error",
			"error", sqlxRow.Err().Error(),
			"order_number", orderNumber)
		return ErrUnexpextedDBError
	}

	var orderEntity dto.OrderEntity
	if err := sqlxRow.StructScan(&orderEntity); err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	_, err = tx.ExecContext(ctx, "select * from user_balance where user_login = $1 for update", orderEntity.UserLogin)
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	_, err = tx.ExecContext(ctx, "update user_balance set total_sum = total_sum + $1 where user_login = $2", accrual, orderEntity.UserLogin)
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	_, err = tx.ExecContext(ctx, "update orders set status = $1, accrual = $2 where number = $3", orderStatus, accrual, orderNumber)
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	err = tx.Commit()
	if err != nil {
		log.Errorw("storage_order: unexpected db error", "error", err.Error())
		return ErrUnexpextedDBError
	}

	return nil
}

func (s *OrderStorage) GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*dto.OrderEntity, error) {
	sqlxRow := s.db.QueryRowxContext(
		ctx,
		"select * from orders where number = $1",
		orderNumber,
	)

	var orderEntity dto.OrderEntity
	if err := sqlxRow.StructScan(&orderEntity); err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrEntityNotFound
	} else if err != nil {
		log.Errorw(
			"storage_user: unexpected error when execute query",
			"error",
			err.Error(),
		)

		return nil, ErrUnexpextedDBError
	}

	return &orderEntity, nil
}

func (s *OrderStorage) GetOrdersByUserLogin(ctx context.Context, userLogin string) ([]dto.OrderEntity, error) {
	rows, err := s.db.QueryxContext(ctx, "select * from orders where user_login = $1", userLogin)
	if err != nil {
		log.Errorw(
			"storage_order: error when get orders by user_login",
			"error", err.Error(),
			"user_login", userLogin,
		)

		return nil, err
	}

	if rows.Err() != nil {
		log.Errorw("storage_order: unexpected DB error", "error", rows.Err().Error())
		return nil, ErrUnexpextedDBError
	}

	var orders = make([]dto.OrderEntity, 0)
	for rows.Next() {
		var order dto.OrderEntity
		if err := rows.StructScan(&order); err != nil {
			log.Errorw(
				"storage_order: error when scan row into dto.OrderEntity",
				"error",
				err.Error(),
			)

			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
