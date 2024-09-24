package dto

import "time"

type APIRegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type APILoginUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderStatus string

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type APIOrderResponse struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type APIUserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type APIPutOrderWithdrawnRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type APIOrderWithdrawn struct {
	OrderNumber   string    `json:"order"`
	WithdrawalSum float64   `json:"sum"`
	ProcessedAt   time.Time `json:"processed_at"`
}
