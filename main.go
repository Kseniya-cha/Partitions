package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/cmd"
	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/pkg/config"
	"github.com/Kseniya-cha/LEARN_GOLANG/partitions/pkg/logger"
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

func covertTime(t time.Time) string {
	tNow := strings.Join(strings.Split(t.Format(time.RFC3339), "T"), " ")
	tNow1 := strings.Split(tNow, ":")

	return strings.Join(tNow1[:len(tNow1)-1], ":")
}

// принимает исходное время t, в которое необходимо запушить в бд
// и число дней, на которое должен распространяться период партиции
func getPeriodDays(t time.Time, period int) (string, string) {
	start, end := covertTime(t), covertTime(t.AddDate(0, 0, period))

	return strings.Split(start, " ")[0] + " 00:00:00+" + strings.Split(start, "+")[1],
		strings.Split(end, " ")[0] + " 00:00:00+" + strings.Split(end, "+")[1]
}

// партиции раз в два часа, только чётные: с 00 до 02, с 12 до 14 etc
func getPeriod2Hour(t time.Time) (string, string) {
	start, end := covertTime(t), covertTime(t.Add(2*time.Hour))

	timezone := strings.Split(start, "+")[1]

	startHour := strings.Split(strings.Split(start, " ")[1], ":")[0]
	endHour := strings.Split(strings.Split(end, " ")[1], ":")[0]

	if i, _ := strconv.Atoi(startHour); i%2 != 0 {
		startHourI, _ := strconv.Atoi(startHour)
		startHour = strconv.Itoa(startHourI - 1)
		endHourI, _ := strconv.Atoi(endHour)
		endHour = strconv.Itoa(endHourI - 1)
	}

	return strings.Split(start, " ")[0] + " " + startHour + ":00:00+" + timezone,
		strings.Split(end, " ")[0] + " " + endHour + ":00:00+" + timezone
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
	// tableName := "my_table"

	// // коннект к бд
	// db, err := NewDB()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	// 	os.Exit(1)
	// }

	// ctx, _ := context.WithCancel(context.Background())

	// // раз в какое время будем пытаться вставить строки
	// tick := time.NewTicker(10 * time.Minute)

	// for {
	// 	<-tick.C

	// 	// время получения сообщения
	// 	t := time.Now()

	// 	// start, end - начало и конец границ партиций
	// 	// start, end := getPeriodDays(t, 1)
	// 	start, end := getPeriod2Hour(t)

	// 	// имя партиции
	// 	partitionName := getPartitionName(tableName, start, end, true)

	// 	// проверка, что партиция существует
	// 	_, err := db.Conn.Query(ctx, fmt.Sprintf(`SELECT * FROM %s`, partitionName))

	// 	// если партиции не существует, она создаётся
	// 	if err != nil {
	// 		_, err = db.Conn.Exec(ctx, fmt.Sprintf(`
	// 			CREATE TABLE %s PARTITION OF %s
	// 			FOR VALUES FROM ('%s') TO ('%s')
	// 		`, partitionName, tableName, start, end))
	// 		if err != nil {
	// 			fmt.Println("cannot create partition table:", err)
	// 			continue
	// 		}
	// 	}

	// 	// формирование запроса
	// 	query := fmt.Sprintf("INSERT INTO %s (name, created_at) VALUES ('test', '%v')", tableName, covertTime(t))

	// 	// вставка строки
	// 	smth, err := db.Conn.Exec(ctx, query)
	// 	fmt.Println(smth, err)
	// 	if err != nil {
	// 		fmt.Println("cannot insert:", err)
	// 	}
	// }

	//
	//
	//

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
