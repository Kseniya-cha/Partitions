package globalcycle

import "time"

type MonitoringCycles struct {
	Id            int       `json: "id"`
	StartDatetime time.Time `json: "start_datetime`
	EndDatetime   time.Time `json: "end_datetime`
}
