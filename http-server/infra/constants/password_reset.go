package constants

import "time"

const (
	PasswordResetSendLimit24h = 3
	PasswordResetTokenTTL     = 24 * time.Hour
	PasswordResetSendLimitTTL = 24 * time.Hour
)
