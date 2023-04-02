package devicetestingresults

const (
	SelectAll = `SELECT * FROM %s`

	CreatePartition = `
	CREATE TABLE %s PARTITION OF %s
	FOR VALUES FROM ('%s') TO ('%s')`

	Insert = `INSERT INTO %s (cycles_id) VALUES ('%d', '%s')`
)
