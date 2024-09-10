package dto

import "time"

type OrderDomain struct {
	Number     string
	Status     OrderStatus
	Accrual    *float64
	UploadedAt time.Time
}

type UserBalanceDomain struct {
	Current  float64
	Withdraw float64
}

type UserWithdrawDomain struct {
	Order       string
	Sum         float64
	ProcessedAt time.Time
}

type AccuralInfoDomain struct {
}
