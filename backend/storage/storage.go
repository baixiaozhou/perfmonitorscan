package storage

import (
	"fmt"
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	MYSQL      = "MYSQL"
	POSTGRESQL = "POSTGRESQL"
	SQLITE     = "SQLITE"
)

func InitDataBase(config *conf.DBConfig) {
	var err error
	dbConfig := conf.GlobalConfig.DB
	switch dbConfig.Database {
	case MYSQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case POSTGRESQL:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.DbName, dbConfig.Port)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case SQLITE:
		DB, err = gorm.Open(sqlite.Open("monitoring.db"), &gorm.Config{})
	default:
		conf.Logger.Fatal("Unsupported database type")
	}
	if err != nil {
		conf.Logger.Fatal("Init DB Error:" + err.Error())
	}

	err = DB.AutoMigrate(&MonitoringData{})
	if err != nil {
		conf.Logger.Fatal("Init DB AutoMigrate Error:" + err.Error())
	}
}

func SaveData(data *MonitoringData) error {
	return DB.Create(&data).Error
}
