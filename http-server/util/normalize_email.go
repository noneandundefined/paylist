package util

import "strings"

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(email), " ", ""))
}
