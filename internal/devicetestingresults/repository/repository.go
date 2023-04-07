package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kseniya-cha/Partitions/internal/devicetestingresults"
	"github.com/Kseniya-cha/Partitions/pkg/database/postgresql"
	"go.uber.org/zap"
)

type repository struct {
	db  postgresql.IDB
	log *zap.Logger

	common devicetestingresults.Common
}

func NewRepository(db postgresql.IDB, log *zap.Logger) *repository {
	return &repository{
		db:  db,
		log: log,
	}
}

func (r *repository) getInsertQuery(objs []devicetestingresults.DeviceTestingResults,
	tableNameResult string) string {

	val := strings.Builder{}

	if len(objs) == 0 {
		return ""
	}

	for _, obj := range objs {
		val.WriteString(fmt.Sprintf("('%d', '%s'), ", obj.CycleId, obj.StartDatetime))
	}

	return fmt.Sprintf(devicetestingresults.Insert, tableNameResult, val.String()[:val.Len()-2])
}

// IsPartitionExist делает select-запрос к бд и проверяет,
// существует ли партиция с переданным именем
func (r *repository) IsPartitionExist(ctx context.Context, partitionName string) error {

	query := fmt.Sprintf(devicetestingresults.SelectAll, partitionName)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	r.log.Debug(fmt.Sprintf("SQL query:\n%s", query))

	tx, _ := r.db.GetConn().Begin(ctx)
	defer tx.Rollback(ctx)
	defer tx.Commit(ctx)

	if _, err := tx.Query(ctx, query); err != nil {
		return err
	}

	return nil
}

// CreatePartition создаёт партицию partitionName к таблице tableName
// с промежутком от start до end
func (r *repository) CreatePartition(ctx context.Context, partitionName, tableName,
	start, end string) error {

	query := fmt.Sprintf(devicetestingresults.CreatePartition, partitionName, tableName, start, end)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	r.log.Debug(fmt.Sprintf("SQL query:\n%s", query))

	tx, _ := r.db.GetConn().Begin(ctx)
	defer tx.Rollback(ctx)
	defer tx.Commit(ctx)

	if _, err := tx.Exec(ctx, query); err != nil {
		return fmt.Errorf("cannot create partition: %v", err)
	}

	return nil
}

func (r *repository) Insert(ctx context.Context, tableNameResult string, objs []devicetestingresults.DeviceTestingResults) error {

	query := r.getInsertQuery(objs, tableNameResult)

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
