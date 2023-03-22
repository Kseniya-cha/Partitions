package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// GracefulShutdown - метод для корректного завершения работы программы
// при получении прерывающего сигнала
func (a *app) GracefulShutdown(cancel context.CancelFunc) {
	defer time.Sleep(time.Second * 5)
	defer close(a.sigChan)
	defer cancel()

	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)

	sign := <-a.sigChan
	a.log.Info(fmt.Sprintf("Got signal: %v, exiting", sign))

	a.db.Close()
	a.log.Info("Close database connection")

	a.log.Debug("Waiting...")
}

func convertTime(t time.Time) string {
	tNow := strings.Join(strings.Split(t.Format(time.RFC3339), "T"), " ")
	tNow1 := strings.Split(tNow, ":")

	return strings.Join(tNow1[:len(tNow1)-1], ":")
}

// принимает исходное время t, в которое необходимо запушить в бд
// и число дней, на которое должен распространяться период партиции
func getPeriodDays(t time.Time, period int) (string, string) {
	start, end := convertTime(t), convertTime(t.AddDate(0, 0, period))

	return strings.Split(start, " ")[0] + " 00:00:00",
		strings.Split(end, " ")[0] + " 00:00:00"
}

// партиции раз в два часа, только чётные: с 00 до 02, с 12 до 14 etc
func getPeriod2Hour(t time.Time) (string, string) {
	start, end := convertTime(t), convertTime(t.Add(2*time.Hour))

	fmt.Println(start)
	fmt.Println(strings.Split(start, " ")[0] + " 00:00:00")

	startHour := strings.Split(strings.Split(start, " ")[1], ":")[0]
	endHour := strings.Split(strings.Split(end, " ")[1], ":")[0]

	if i, _ := strconv.Atoi(startHour); i%2 != 0 {
		startHourI, _ := strconv.Atoi(startHour)
		startHour = strconv.Itoa(startHourI - 1)
		endHourI, _ := strconv.Atoi(endHour)
		endHour = strconv.Itoa(endHourI - 1)
	}

	return strings.Split(start, " ")[0] + " " + startHour + ":00:00",
		strings.Split(end, " ")[0] + " " + endHour + ":00:00"
}

func getPartitionName(tableName, start, end string, isHour bool) string {
	partitionName := tableName + "_" + strings.Replace(strings.Split(start, " ")[0], "-", "_", -1)

	if isHour {
		hourStart := strings.Split(strings.Split(start, " ")[1], ":")
		hourEnd := strings.Split(strings.Split(end, " ")[1], ":")
		partitionName += "_from" + hourStart[0] + "to" + hourEnd[0]
	}

	return partitionName
}
