package repository

import (
	"context"
	"fmt"

	"github.com/Kseniya-cha/Partitions/internal/monitoringcycle"
	"github.com/Kseniya-cha/Partitions/pkg/database/postgresql"
	"go.uber.org/zap"
)

type repository struct {
	db  postgresql.IDB
	log *zap.Logger

	common monitoringcycle.Common
}

func NewRepository(db postgresql.IDB, log *zap.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (r *repository) InsertGlobal(ctx context.Context, tableName string) error {
	query := fmt.Sprintf(monitoringcycle.InsertGlobal, tableName)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	r.log.Debug(fmt.Sprintf("SQL query:\n%s", query))

	tx, _ := r.db.GetConn().Begin(ctx)
	defer tx.Rollback(ctx)
	defer tx.Commit(ctx)

	if _, err := tx.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

// GetNewGlobalCycle возвращает последний id из таблицы tableName
func (r *repository) GetNewGlobalCycle(ctx context.Context, tableName string) (int, error) {

	query := fmt.Sprintf(monitoringcycle.GetNewGlobalCycle, tableName)

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	r.log.Debug(fmt.Sprintf("SQL query:\n%s", query))

	tx, _ := r.db.GetConn().Begin(ctx)
	defer tx.Rollback(ctx)
	defer tx.Commit(ctx)

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *repository) Update(ctx context.Context, tableName string, id int) error {
	query := fmt.Sprintf(monitoringcycle.Update, tableName, id)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	r.log.Debug(fmt.Sprintf("SQL query:\n%s", query))

	tx, _ := r.db.GetConn().Begin(ctx)
	defer tx.Rollback(ctx)
	defer tx.Commit(ctx)

	if _, err := tx.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}
