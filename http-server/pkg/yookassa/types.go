package yookassa

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type CreatePaymentRequest struct {
	Amount            Amount            `json:"amount"`
	Capture           bool              `json:"capture"`
	Confirmation      Confirmation      `json:"confirmation"`
	Description       string            `json:"description"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	SavePaymentMethod bool              `json:"save_payment_method,omitempty"`
}

type PaymentMethod struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Saved bool   `json:"saved"`
}

type PaymentConfirmation struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type Payment struct {
	ID            string              `json:"id"`
	Status        string              `json:"status"`
	Paid          bool                `json:"paid"`
	Amount        Amount              `json:"amount"`
	Description   string              `json:"description"`
	Metadata      map[string]string   `json:"metadata"`
	PaymentMethod PaymentMethod       `json:"payment_method"`
	Confirmation  PaymentConfirmation `json:"confirmation"`
}

type WebhookNotification struct {
	Type   string  `json:"type"`
	Event  string  `json:"event"`
	Object Payment `json:"object"`
}
