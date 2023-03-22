package cmd

import (
	"context"
	"fmt"
	"time"
)

func (a *app) Run(ctx context.Context) {

	query := fmt.Sprintf("INSERT INTO monitoring_cycles (start_datetime,end_datetime) VALUES (now(), now() + interval '1 second')")

	time.Sleep(time.Second)
	// вставка строки
	_, err := a.db.GetConn().Exec(ctx, query)
	if err != nil {
		return
	}

	// раз в какое время будем пытаться вставить строки
	tick := time.NewTicker(5 * time.Second)

	for {
		<-tick.C

		// время получения сообщения
		t := time.Now()

		// start, end - начало и конец границ партиций
		start, end := getPeriodDays(t, 1)
		// start, end := getPeriod2Hour(t)

		// имя партиции
		partitionName := getPartitionName(a.cfg.TableName, start, end, false)

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
		query := fmt.Sprintf("INSERT INTO %s (cycles_id) VALUES ('1')", a.cfg.TableName)

		// вставка строки
		smth, err := a.db.GetConn().Exec(ctx, query)
		fmt.Println(smth, err)
		if err != nil {
			fmt.Println("cannot insert:", err)
		}

	}
}
