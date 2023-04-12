package methods

import (
	"strings"
	"time"
)

// convertTime преобразует время в строку
func ConvertTime(t time.Time) string {
	tNow := strings.Join(strings.Split(t.Format(time.RFC3339), "T"), " ")
	tNow1 := strings.Split(tNow, ":")

	return strings.Join(tNow1[:len(tNow1)-1], ":")
}
