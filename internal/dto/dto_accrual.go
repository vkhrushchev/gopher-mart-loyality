package dto

type AccrualStatus string

var (
	AccrualRegistredStatus  AccrualStatus = "REGISTERED"
	AccrualInvalidStatus    AccrualStatus = "INVALID"
	AccrualProcessingStatus AccrualStatus = "PROCESSING"
	AccrualProcessedStatus  AccrualStatus = "PROCESSED"
)

type AccuralOrderResponse struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual float64       `json:"accrual"`
}
