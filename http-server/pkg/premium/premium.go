package premium

import (
	"strings"

	"paylist.server/types"
)

func IsPremiumPlan(authToken *types.AuthToken) bool {
	return strings.EqualFold(strings.TrimSpace(authToken.PlanName), "Premium")
}
