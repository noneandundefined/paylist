package constants

import "time"

const (
	EmailConfirmSendLimit24h = 3
	EmailConfirmPendingTTL   = 24 * time.Hour
	EmailConfirmSendLimitTTL = 24 * time.Hour
)
