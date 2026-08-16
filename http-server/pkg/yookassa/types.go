package yookassa

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type,omitempty"`
	ReturnURL string `json:"return_url,omitempty"`
}

type CreatePaymentRequest struct {
	Amount            Amount            `json:"amount"`
	Capture           bool              `json:"capture"`
	Confirmation      *Confirmation     `json:"confirmation,omitempty"`
	PaymentMethodID   string            `json:"payment_method_id,omitempty"`
	Description       string            `json:"description"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	SavePaymentMethod bool              `json:"save_payment_method,omitempty"`
}

type PaymentMethodCard struct {
	First6        string `json:"first6"`
	Last4         string `json:"last4"`
	ExpiryYear    string `json:"expiry_year"`
	ExpiryMonth   string `json:"expiry_month"`
	CardType      string `json:"card_type"`
	IssuerCountry string `json:"issuer_country"`
	IssuerName    string `json:"issuer_name"`
}

type PaymentMethod struct {
	ID    string             `json:"id"`
	Type  string             `json:"type"`
	Title string             `json:"title"`
	Saved bool               `json:"saved"`
	Card  *PaymentMethodCard `json:"card,omitempty"`
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
