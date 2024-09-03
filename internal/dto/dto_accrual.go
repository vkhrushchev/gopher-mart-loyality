package dto

type AccrualStatus string

type AccuralOrderResponse struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual float64       `json:"accrual"`
}
