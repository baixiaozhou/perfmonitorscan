package test

import (
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"github.com/baixiaozhou/perfmonitorscan/backend/storage"
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

func TestES(t *testing.T) {
	err := conf.LoadConfig("../config.yml")
	if err != nil {
		t.Error(err)
	}
	err = storage.InitDataBase(&conf.GlobalConfig.DB)
	if err != nil {
		t.Error(err)
	}
	//storage.getDataNodeCount()
}
