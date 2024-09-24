package dto

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
}

type OrderEntity struct {
	ID         int64           `db:"id"`
	UserLogin  string          `db:"user_login"`
	Number     string          `db:"number"`
	Status     OrderStatus     `db:"status"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UploadedAt time.Time       `db:"uploaded_at"`
}

type UserBalanceEntity struct {
	ID                 int64   `db:"id"`
	UserLogin          string  `db:"user_login"`
	TotalSum           float64 `db:"total_sum"`
	TotalWithdrawalSum float64 `db:"total_withdrawal_sum"`
}

type BalanceWithdrawalEntity struct {
	ID            int64     `db:"id"`
	UserLogin     string    `db:"user_login"`
	OrderNumber   string    `db:"order_number"`
	WithdrawalSum float64   `db:"withdrawal"`
	ProcessedAt   time.Time `db:"processed_at"`
}
