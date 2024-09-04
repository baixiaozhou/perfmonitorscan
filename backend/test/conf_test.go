package test

import (
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"testing"
)

func TestConf(t *testing.T) {
	err := conf.LoadConfig("../config.yml")
	if err != nil {
		t.Error(err)
	}

	exceptDB := "postgresql"
	if exceptDB != conf.GlobalConfig.DB.Database {
		t.Errorf("Unsupport")
	}
}
