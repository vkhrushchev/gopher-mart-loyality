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
