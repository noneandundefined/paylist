package payment_handler_v1

type AutoRenewPayload struct {
	Enabled bool `json:"enabled"`
}

type CheckoutPayload struct {
	PlanName string `json:"plan_name" validate:"required,min=2,max=64"`
}

type CheckoutResponse struct {
	PaymentID       string `json:"payment_id"`
	ConfirmationURL string `json:"confirmation_url"`
}

type PaymentConfirmResponse struct {
	Paid   bool   `json:"paid"`
	Status string `json:"status"`
}

type PaymentHistoryResponse struct {
	ID          uint64  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	PlanName    string  `json:"plan_name"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	PaymentKind string  `json:"payment_kind"`
	Description *string `json:"description,omitempty"`
	PaidAt      *string `json:"paid_at,omitempty"`
}
