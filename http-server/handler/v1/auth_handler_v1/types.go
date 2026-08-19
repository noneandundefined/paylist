package auth_handler_v1

type AuthSignupPayload struct {
	FirstName    *string `json:"first_name" validate:"omitempty,min=3,max=45"`
	LastName     *string `json:"last_name" validate:"omitempty,min=3,max=45"`
	Email        string  `json:"email" validate:"required,email"`
	Password     string  `json:"password" validate:"required,min=6,max=16"`
	ReferralCode string  `json:"referral_code" validate:"omitempty,max=16"`
}

type AuthSigninPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=6,max=16"`
}

type AuthCheckPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqEmailConfirmPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthPasswordResetRequestPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthPasswordResetConfirmPayload struct {
	Password string `json:"password" validate:"required,min=6,max=16"`
}
