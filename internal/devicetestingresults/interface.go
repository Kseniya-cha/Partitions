package devicetestingresults

import (
	"context"
	"time"
)

type Repository interface {
	IsPartitionExist(ctx context.Context, partitionName string) error
	CreatePartition(ctx context.Context, partitionName, tableName, start, end string) error
	Insert(ctx context.Context, tableNameResult string,
		objs []DeviceTestingResults, t time.Time) error
}

type Common interface {
	Repository
}
