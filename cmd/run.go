package cmd

import (
	"context"
	"math/rand"
	"time"

	"github.com/Kseniya-cha/Partitions/internal/devicetestingresults"
	"github.com/Kseniya-cha/Partitions/pkg/methods"
)

func (a *app) Run(ctx context.Context) {

	// раз в какое время будем пытаться вставить строки
	tick := time.NewTicker(5 * time.Second)

	for {
		<-tick.C

		if ctx.Err() != nil {
			return
		}

		if !a.db.IsConn(ctx) {
			return
		}

		err := a.repoGlobal.InsertGlobal(ctx, a.cfg.TableNameGlobal)
		if err != nil {
			a.log.Error(err.Error())
			continue
		}

		// // start, end - начало и конец границ партиций
		// start, end := getPeriodDays(t, 1)

		// // имя партиции
		// partitionName := getPartitionName(a.cfg.TableNameResult, start, end, false)

		// // проверка, что партиция существует
		// err = a.repoResult.IsPartitionExist(ctx, partitionName)

		// // если партиции не существует, она создаётся
		// if err != nil {
		// 	err = a.repoResult.CreatePartition(ctx, partitionName, a.cfg.TableNameResult, start, end)
		// 	if err != nil {
		// 		a.log.Error(err.Error())
		// 		continue
		// 	}
		// }

		// получений id нового цикла
		id, err := a.repoGlobal.GetNewGlobalCycle(ctx, a.cfg.TableNameGlobal)
		if err != nil {
			a.log.Error(err.Error())
			continue
		}

		// объекты, которые необходимо вставить за данный цикл
		objs := []devicetestingresults.DeviceTestingResults{
			{CycleId: id, StartDatetime: methods.ConvertTime(time.Now()), Uuid: rand.Intn(1000)},
			{CycleId: id, StartDatetime: methods.ConvertTime(time.Now()), Uuid: rand.Intn(1000)},
			{CycleId: id, StartDatetime: methods.ConvertTime(time.Now()), Uuid: rand.Intn(1000)},
		}

		t := time.Now()

		// вставка строки в таблицу с результатами
		err = a.repoResult.Insert(ctx, a.cfg.TableNameResult, objs, t)
		if err != nil {
			a.log.Error(err.Error())
			continue
		}

		// вставка конечного времени обработки цикла
		err = a.repoGlobal.Update(ctx, a.cfg.TableNameGlobal, id)
		if err != nil {
			a.log.Error(err.Error())
		}
	}
}
