package repository

import (
	"strings"
	"time"

	"github.com/Kseniya-cha/Partitions/pkg/methods"
)

// getPeriodDays принимает исходное время t, в которое необходимо запушить в бд
// и число дней, на которое должен распространяться период партиции
func getPeriodDays(t time.Time, period int) (string, string) {
	start, end := methods.ConvertTime(t), methods.ConvertTime(t.AddDate(0, 0, period))

	return strings.Split(start, " ")[0] + " 00:00:00",
		strings.Split(end, " ")[0] + " 00:00:00"
}

func getPartitionName(tableName, start string) string {
	return tableName + "_" + strings.Replace(strings.Split(start, " ")[0], "-", "_", -1)
}
