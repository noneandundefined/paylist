package device_handler_v1

type DeviceCreateAuthSessionPayload struct {
	DeviceId string `json:"device_id" validate:"required,min=5"`
}

type DeviceConfirmPayload struct {
	SessionId string `json:"session_id" validate:"required,min=5"`
}
