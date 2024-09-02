package api

import (
	"fmt"
	"github.com/baixiaozhou/perfmonitorscan/backend/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CollectCpuData(c *gin.Context) {
	var data storage.MonitoringCpuData

	// test
	fmt.Println("receive data")
	// test
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// test
	fmt.Println("receive data, data is:", data)
	// test
	//if err := storage.SaveData(&data); err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//}

	c.JSON(http.StatusOK, gin.H{"status": "Data collected successfully"})
}
