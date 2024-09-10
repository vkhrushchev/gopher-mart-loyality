package dto

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	Id           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
}

type OrderEntity struct {
	Id         int64           `db:"id"`
	UserLogin  string          `db:"user_login"`
	Number     string          `db:"number"`
	Status     OrderStatus     `db:"status"`
	Accrual    sql.NullFloat64 `db:"accrual"`
	UploadedAt time.Time       `db:"uploaded_at"`
}
