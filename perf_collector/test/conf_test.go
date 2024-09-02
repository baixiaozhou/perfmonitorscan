package test

import (
	"github.com/baixiaozhou/perfmonitorscan/perf_collector/collector"
	"github.com/baixiaozhou/perfmonitorscan/perf_collector/conf"
	"testing"
)

func TestParseConf(t *testing.T) {
	err := conf.LoadConfig("../config.yml")
	if err != nil {
		t.Error(err)
	}

	conf.InitLogger(&conf.GlobalConfig.Log)
	logger := conf.GetLogger()
	logger.Info("hello")
}

func TestJStackAnalysis(t *testing.T) {
	info, _err := collector.GenerJstackInfo("jstack.txt")
	if _err != nil {
		t.Fatal(_err)
	}
	expectedDeadLockCount := 2
	if info.DeadLockCount != expectedDeadLockCount {
		t.Errorf("expectedDeadLockCount is %d, but got %d", expectedDeadLockCount, info.DeadLockCount)
	}
	expectedThreadCount := 73
	if info.ThreadsCount != expectedThreadCount {
		t.Errorf("expectedThreadCount is %d, but got %d", expectedThreadCount, info.ThreadsCount)
	}
	expectedRunnableCount := 22
	if info.RunnableCount != expectedRunnableCount {
		t.Errorf("expectedRunnableCount is %d, but got %d", expectedRunnableCount, info.RunnableCount)
	}
}
