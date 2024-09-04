package api

import (
	"github.com/baixiaozhou/perfmonitorscan/backend/conf"
	"github.com/baixiaozhou/perfmonitorscan/backend/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CollectCpuData(c *gin.Context) {
	var data storage.MonitoringCpuData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := storage.SaveData(&data); err != nil {
		conf.Logger.Error("failed to save data to db:", err.Error(), " data:", data)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": "Data collected successfully"})
}
