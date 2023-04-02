package monitoringcycle

import "context"

type Repository interface {
	InsertGlobal(ctx context.Context, tableName string) error
	GetNewGlobalCycle(ctx context.Context, tableName string) (int, error)
	Update(ctx context.Context, tableName string, id int) error
}

type Common interface {
	Repository
}
