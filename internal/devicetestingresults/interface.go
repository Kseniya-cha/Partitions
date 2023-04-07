package devicetestingresults

import "context"

type Repository interface {
	IsPartitionExist(ctx context.Context, partitionName string) error
	CreatePartition(ctx context.Context, partitionName, tableName, start, end string) error
	Insert(ctx context.Context, tableNameResult string, objs []DeviceTestingResults) error

	getInsertQuery(objs []DeviceTestingResults, tableNameResult string) string
}

type Common interface {
	Repository
}
