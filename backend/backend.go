package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/baixiaozhou/perfmonitorscan/backend/api"
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"github.com/baixiaozhou/perfmonitorscan/backend/storage"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {
	listenPort := kingpin.Flag(
		"listen-port",
		"Port to listen on.",
	).String()
	config_path := kingpin.Flag(
		"config-path",
		"Path to the configuration file.",
	).String()
	refresh_interval := kingpin.Flag(
		"refresh-interval",
		"Refresh interval in seconds.",
	).Default("60").Int()
	kingpin.Parse()

	err := conf.LoadConfig(*config_path)
	if err != nil {
		log.Fatal(err)
	}

	conf.InitLogger(&conf.GlobalConfig.Logging)

	ticker := time.NewTicker(time.Duration(*refresh_interval) * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				err := conf.LoadConfig(*config_path)
				if err != nil {
					log.Fatal(err)
				} else {
					conf.ReloadLogConfig(&conf.GlobalConfig.Logging)
				}
			}
		}
	}()

	if err := storage.InitDataBase(&conf.GlobalConfig.DB); err != nil {
		fmt.Println("err:", err)
		conf.Logger.Fatal(err)
	}

	conf.Logger.Info("listening on port " + *listenPort)
	// Create gin
	r := gin.Default()
	r.POST("/collect/cpu", api.CollectCpuData)
	r.Run(*listenPort)
}
