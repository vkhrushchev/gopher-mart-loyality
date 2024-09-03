package dto

import "time"

type APIRegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type APILoginUserReqest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OrderStatus string

type APIGetUserOrderResponseEntry struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

type APIGetUserBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type APIWithdrawUserBalanceRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type APIGetUserBalanaceWithdrawlsResponseEntry struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
