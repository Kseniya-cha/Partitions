package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Conn *pgxpool.Pool
}

func NewDB() (db *DB, err error) {

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}
	config.ConnConfig.User = "postgres"
	config.ConnConfig.Password = "password"
	config.ConnConfig.Host = "localhost"
	config.ConnConfig.Port = uint16(5432)
	config.ConnConfig.Database = "firstbd"

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	return &DB{pool}, nil
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

func main() {
	tableName := "device_testing_results"

	// коннект к бд
	db, err := NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	ctx, _ := context.WithCancel(context.Background())

	query := fmt.Sprintf("INSERT INTO monitoring_cycles (start_datetime,end_datetime) VALUES (now(), now() + interval '1 second')")

	time.Sleep(time.Second)
	// вставка строки
	_, err = db.Conn.Exec(ctx, query)

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
		partitionName := getPartitionName(tableName, start, end, false)

		// проверка, что партиция существует
		_, err := db.Conn.Query(ctx, fmt.Sprintf(`SELECT * FROM %s`, partitionName))

		// если партиции не существует, она создаётся
		if err != nil {
			_, err = db.Conn.Exec(ctx, fmt.Sprintf(`
				CREATE TABLE %s PARTITION OF %s
				FOR VALUES FROM ('%s') TO ('%s')
			`, partitionName, tableName, start, end))
			if err != nil {
				fmt.Println("cannot create partition table:", err)
				continue
			}
		}

		// формирование запроса
		query := fmt.Sprintf("INSERT INTO %s (cycles_id) VALUES ('1')", tableName)

		// вставка строки
		smth, err := db.Conn.Exec(ctx, query)
		fmt.Println(smth, err)
		if err != nil {
			fmt.Println("cannot insert:", err)
		}

	}
}
