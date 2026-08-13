package categoryslug

import (
	"strings"
	"unicode"
)

func FromLabel(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		return ""
	}

	var builder strings.Builder
	prevDash := false

	for _, r := range strings.ToLower(label) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			prevDash = false
			continue
		}

		if !prevDash {
			builder.WriteRune('-')
			prevDash = true
		}
	}

	slug := strings.Trim(builder.String(), "-")
	if len(slug) > 64 {
		slug = slug[:64]
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}
