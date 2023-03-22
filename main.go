package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/cmd"
	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/pkg/config"
	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/pkg/logger"
)

func main() {
	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())

	// Чтение конфигурационного файла
	cfg, err := config.GetConfig()
	if err != nil || reflect.DeepEqual(cfg.Database, config.Database{}) {
		fmt.Println("ERROR: cannot read config: file is empty")
		return
	}

	log := logger.NewLogger(cfg)

	// Инициализация прототипа приложения
	app, err := cmd.NewApp(ctx, cfg)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Запуск алгоритма в отдельной горутине
	go app.Run(ctx)

	// Ожидание прерывающего сигнала
	app.GracefulShutdown(cancel)
}
