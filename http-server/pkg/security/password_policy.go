package security

import "strings"

func PasswordEqualsEmail(password, email string) bool {
	password = strings.ToLower(strings.TrimSpace(password))
	email = strings.ToLower(strings.TrimSpace(email))

	if password == "" || email == "" {
		return false
	}

	if password == email {
		return true
	}

	if at := strings.Index(email, "@"); at > 0 {
		if password == email[:at] {
			return true
		}
	}

	return false
}
