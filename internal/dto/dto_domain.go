package dto

import "time"

type OrderDomain struct {
	Number     string
	Status     OrderStatus
	Accrual    *float64
	UploadedAt time.Time
}

type UserBalanceDomain struct {
	Current    float64
	Withdrawal float64
}

type OrderWithdrawalDomain struct {
	OrderNumber   string
	WithdrawalSum float64
	ProcessedAt   time.Time
}

type AccuralInfoDomain struct {
	OrderNumber string
	OrderStatus AccrualStatus
	Accrual     float64
}
