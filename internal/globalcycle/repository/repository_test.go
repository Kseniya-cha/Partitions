package repository

import (
	"testing"
	"time"

	"github.com/Kseniya-cha/Partitions/pkg/methods"
)

func Test_convertTime(t *testing.T) {
	expect := "2001-01-01 12:00:00+03"
	got := methods.ConvertTime(time.Date(2001, time.January, 1, 12, 0, 0, 0, time.Local))
	if got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
}

func Test_getPeriodDays(t *testing.T) {

	startTime := time.Date(2001, time.January, 1, 12, 0, 0, 0, time.Local)

	startExp, endExp := "2001-01-01 00:00:00", "2001-01-02 00:00:00"
	startGot, endGot := getPeriodDays(startTime, 1)

	if startExp != startGot || endExp != endGot {
		t.Errorf("expect: %s, got: %s", startExp, startGot)
	}

	if endExp != endGot {
		t.Errorf("expect: %s, got: %s", endExp, endGot)
	}
}

func Test_getPartitionName(t *testing.T) {
	tableName, start := "test", "2001-01-01 00:00:00"
	expect := "test_2001_01_01"
	got := getPartitionName(tableName, start)

	if got != expect {
		t.Errorf("expect: %s, got: %s", expect, got)
	}
}
