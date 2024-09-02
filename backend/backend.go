package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"github.com/baixiaozhou/perfmonitorscan/backend/controller"
	"github.com/baixiaozhou/perfmonitorscan/backend/storage"
	"github.com/gin-gonic/gin"
	"log"
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
	kingpin.Parse()

	err := conf.LoadConfig(*config_path)
	if err != nil {
		log.Fatal(err)
	}

	conf.InitLogger(&conf.GlobalConfig.Logging)

	storage.InitDataBase(&conf.GlobalConfig.DB)
	// Create gin
	r := gin.Default()
	r.POST("/collect", controller.CollectData)
	r.Run(*listenPort)
}
