package cmd

import (
	"context"
	"fmt"
	"os/signal"
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

// convertTime преобразует время в строку
func convertTime(t time.Time) string {
	tNow := strings.Join(strings.Split(t.Format(time.RFC3339), "T"), " ")
	tNow1 := strings.Split(tNow, ":")

	return strings.Join(tNow1[:len(tNow1)-1], ":")
}
