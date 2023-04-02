package devicetestingresults

type DeviceTestingResults struct {
	Id            int    `json: "id"`
	CycleId       int    `json: "cycles_id"`
	StartDatetime string `json: "start_datetime"`
}
