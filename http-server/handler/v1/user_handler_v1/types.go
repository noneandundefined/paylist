package user_handler_v1

type UserSessionResponse struct {
	SessionId string `json:"session_id"`
	Platform  string `json:"platform"`
	DeviceId  string `json:"device_id,omitempty"`
	CreatedAt string `json:"created_at"`
	Current   bool   `json:"current"`
}

type UserSessionDisconnectPayload struct {
	SessionId string `json:"session_id" validate:"required,min=5"`
}

type UserProfileUpdatePayload struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=3,max=45"`
	LastName  *string `json:"last_name" validate:"omitempty,min=3,max=45"`
}

type UserEmailChangePayload struct {
	Email string `json:"email" validate:"required,email"`
}

type UserSettingsResponse struct {
	DisplayCurrency   *string `json:"display_currency,omitempty"`
	Country           *string `json:"country,omitempty"`
	TelegramConnected bool    `json:"telegram_connected"`
	TelegramUsername  *string `json:"telegram_username,omitempty"`
}

type UserSettingsUpdatePayload struct {
	DisplayCurrency *string `json:"display_currency" validate:"omitempty,len=3"`
	Country         *string `json:"country" validate:"omitempty,len=2"`
}
