package globalcycle

const (
	InsertGlobal      = `INSERT INTO %s (start_datetime) VALUES (now())`
	GetNewGlobalCycle = `SELECT id FROM %s ORDER BY "id" DESC LIMIT 1`
	Update            = `UPDATE %s SET end_datetime=now() WHERE id=%d`
)
