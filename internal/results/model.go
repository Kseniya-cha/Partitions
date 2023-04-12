package results

type DeviceTestingResults struct {
	Id            int    `json: "id"`
	CycleId       int    `json: "cycles_id"`
	Uuid          int    `json: "uuid"`
	StartDatetime string `json: "start_datetime"`
}
