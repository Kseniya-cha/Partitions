package main

import (
	"context"
	"fmt"
	"os"
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

func covertTime(t time.Time) string {
	tNow := strings.Join(strings.Split(t.Format(time.RFC3339), "T"), " ")
	tNow1 := strings.Split(tNow, ":")

	return strings.Join(tNow1[:len(tNow1)-1], ":")
}

func getStartEndTime(t time.Time) (string, string) {
	start, end := covertTime(t), covertTime(t.AddDate(0, 0, 1))

	return strings.Split(start, " ")[0] + " 00:00:00+" + strings.Split(start, "+")[1],
		strings.Split(end, " ")[0] + " 00:00:00+" + strings.Split(end, "+")[1]
}

func main() {
	tableName := "my_table"

	db, err := NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	ctx, _ := context.WithCancel(context.Background())

	for {
		time.Sleep(5 * time.Second)

		t := time.Now()
		start, end := getStartEndTime(t)

		partitionName := tableName + "_" + strings.Replace(strings.Split(covertTime(t), " ")[0], "-", "_", -1)

		_, err := db.Conn.Query(ctx, fmt.Sprintf(`SELECT * FROM %s`, partitionName))

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

		query := fmt.Sprintf("INSERT INTO %s (name, created_at) VALUES ('test', '%v')", tableName, covertTime(t))

		smth, err := db.Conn.Exec(ctx, query)
		fmt.Println(smth, err)
		if err != nil {
			fmt.Println("cannot insert:", err)
		}
	}
}
