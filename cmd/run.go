package cmd

import (
	"context"
	"math/rand"
	"time"

	"github.com/Kseniya-cha/Partitions/internal/results"
	"github.com/Kseniya-cha/Partitions/pkg/methods"
)

func (a *app) Run(ctx context.Context) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// период запуска цикла
	tick := time.NewTicker(5 * time.Second)

	for {
		<-tick.C

		if ctx.Err() != nil {
			return
		}

		// добавление информации о новом цикле
		err := a.repoGlobal.InsertGlobal(ctx, a.cfg.TableNameGlobal)
		if err != nil {
			a.log.Error(err.Error())
			continue
		}

		// получение id нового цикла
		id, err := a.repoGlobal.GetNewGlobalCycle(ctx, a.cfg.TableNameGlobal)
		if err != nil {
			a.log.Error(err.Error())
			continue
		}

		t := time.Now()

		// рандомное число объектов, которые необходимо вставить за данный цикл
		objs := []results.DeviceTestingResults{}
		for i := 0; i < rand.Intn(5); i++ {
			objs = append(objs, results.DeviceTestingResults{CycleId: id, StartDatetime: methods.ConvertTime(t), Uuid: rand.Intn(1000)})
		}

		// вставка строки в таблицу с результатами
		err = a.repoResult.Insert(ctx, a.cfg.TableNameResult, objs, t)
		if err != nil {
			a.log.Error(err.Error())
		}

		// вставка конечного времени обработки цикла
		err = a.repoGlobal.Update(ctx, a.cfg.TableNameGlobal, id)
		if err != nil {
			a.log.Error(err.Error())
		}
	}
}
