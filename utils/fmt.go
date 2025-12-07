package utils

import (
	"fmt"
	"time"
)

func FormatDuration(seconds int64) string {
	if seconds <= 0 {
		return "0s"
	}

	d := time.Duration(seconds) * time.Second

	h := int64(d.Hours())
	m := int64(d.Minutes()) % 60
	s := int64(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
