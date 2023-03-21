package cmd

import (
	"context"
	"fmt"
	"time"
)

func (a *app) Run(ctx context.Context) {

	// раз в какое время будем пытаться вставить строки
	tick := time.NewTicker(10 * time.Minute)

	for {
		<-tick.C

		// время получения сообщения
		t := time.Now()

		// start, end - начало и конец границ партиций
		// start, end := getPeriodDays(t, 1)
		start, end := getPeriod2Hour(t)

		// имя партиции
		partitionName := getPartitionName(a.cfg.TableName, start, end, true)

		// проверка, что партиция существует
		_, err := a.db.GetConn().Query(ctx, fmt.Sprintf(`SELECT * FROM %s`, partitionName))

		// если партиции не существует, она создаётся
		if err != nil {
			_, err = a.db.GetConn().Exec(ctx, fmt.Sprintf(`
				CREATE TABLE %s PARTITION OF %s
				FOR VALUES FROM ('%s') TO ('%s')
			`, partitionName, a.cfg.TableName, start, end))
			if err != nil {
				fmt.Println("cannot create partition table:", err)
				continue
			}
		}

		// формирование запроса
		query := fmt.Sprintf("INSERT INTO %s (name, created_at) VALUES ('test', '%v')", a.cfg.TableName, covertTime(t))

		// вставка строки
		smth, err := a.db.GetConn().Exec(ctx, query)
		fmt.Println(smth, err)
		if err != nil {
			fmt.Println("cannot insert:", err)
		}
	}
}
