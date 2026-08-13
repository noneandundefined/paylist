package date

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(d.Time.Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.TrimSpace(strings.Trim(string(b), `"`))
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	layouts := []string{time.DateOnly, time.RFC3339, "2006-01-02T15:04:05Z07:00"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			d.Time = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
			return nil
		}
	}

	return fmt.Errorf("invalid date %q, expected YYYY-MM-DD", s)
}
