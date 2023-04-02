package devicetestingresults

import "context"

type Repository interface {
	IsPartitionExist(ctx context.Context, partitionName string) error
	CreatePartition(ctx context.Context, partitionName, tableName, start, end string) error
	Insert(ctx context.Context, tableName string, values DeviceTestingResults) error
}

type Common interface {
	Repository
}
