package repository

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/Kseniya-cha/Partitions/internal/results"
	"github.com/Kseniya-cha/Partitions/pkg/database/postgresql"
	"github.com/Kseniya-cha/Partitions/pkg/methods"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
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

func Test_NewRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := postgresql.NewMockIDB(ctrl)
	mockDB.EXPECT().Close()
	defer mockDB.Close()
	mockLog := zap.NewNop()

	repo := NewRepository(mockDB, mockLog)
	repoS := strings.Split(fmt.Sprint(repo), " ")
	testRepoS := strings.Split(fmt.Sprint(&repository{db: mockDB, log: mockLog}), " ")

	for i := range repoS {
		if repoS[i] != testRepoS[i] {
			t.Errorf("Unexpected Repository struct: %v, expect %v", testRepoS, repoS)
		}
	}
}

func Test_Insert(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := postgresql.NewMockIDB(ctrl)
	mockDB.EXPECT().Close()
	defer mockDB.Close()
	mockLog := zap.NewNop()
	mockCommon := results.NewMockCommon(ctrl)

	repo := NewRepository(mockDB, mockLog)
	repo.common = mockCommon

	tnow := time.Now()
	objsTest := []results.DeviceTestingResults{
		{CycleId: 1, StartDatetime: methods.ConvertTime(tnow), Uuid: rand.Intn(1000)},
	}

	mockCommon.EXPECT().Insert(ctx, "results", objsTest, tnow)

	t.Run("Insert_OK", func(t *testing.T) {
		if err := repo.common.Insert(ctx, "results", objsTest, tnow); err != nil {
			t.Error("unexpected error:", err)
		}
	})
}
