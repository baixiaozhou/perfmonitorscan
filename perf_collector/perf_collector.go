package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/baixiaozhou/backend/collector"
	"github.com/baixiaozhou/backend/conf"
	"log"
	"time"
)

// var config *conf.PerfCollectorConf

func main() {
	var (
		config_path = kingpin.Flag(
			"config-path",
			"Path to the configuration file.",
		).String()
		//listen_port = kingpin.Flag(
		//	"listen-port",
		//	"Port to listen on.",
		//).Default("6660").Int()
		refresh_interval = kingpin.Flag(
			"refresh-interval",
			"Refresh interval in seconds.",
		).Default("60").Int()
		worker_threads = kingpin.Flag(
			"worker-threads",
			"Number of worker threads.",
		).Default("5").Int()
	)
	kingpin.Parse()
	err := conf.LoadConfig(*config_path)
	if err != nil {
		log.Fatal(err)
	}

	conf.InitLogger(&conf.GlobalConfig.Log)
	logger := conf.GetLogger()
	// API listen
	jsonData, _ := json.Marshal(conf.GlobalConfig)
	fmt.Println(string(jsonData))
	// refresh config
	ticker := time.NewTicker(time.Duration(*refresh_interval) * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				err := conf.LoadConfig(*config_path)
				if err != nil {
					log.Fatal(err)
				}
				conf.ReloadLogConfig(&conf.GlobalConfig.Log)
			}
		}
	}()

	logger.Info("per collector started")

	go collector.Monitor(conf.GlobalConfig.ProcessMonitor, *worker_threads)

	select {}
}
