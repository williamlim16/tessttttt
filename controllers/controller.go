package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"trash-separator/structs"

	"github.com/gin-gonic/gin"
)

func (idb *InDB) GetPlasticSpaceCount(c *gin.Context) {

}

func (idb *InDB) GetMetalSpaceCount(c *gin.Context) {

}

func (idb *InDB) GetGlassSpaceCount(c *gin.Context) {

}

func (idb *InDB) SendLog(c *gin.Context) {
	var (
		result gin.H
		status string
		msg    string
	)

	trashType := c.PostForm("type")
	trashcanID, _ := strconv.Atoi(c.PostForm("trash_can_id"))

	fmt.Println("Trash Type:", trashType)

	insertLog := structs.Logs{
		Type:         trashType,
		Trash_can_id: trashcanID,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
	}

	resultInsertLog := idb.DB.Table("logs").Create(&insertLog)

	if resultInsertLog.Error == nil {
		status = "success"
		msg = "Log successfully added"

		result = gin.H{
			"type":         trashType,
			"trash_can_id": trashcanID,
			"status":       status,
			"message":      msg,
		}

		c.JSON(http.StatusOK, result)

	} else {
		status = "error"
		msg = "Log insertion failed"

		result = gin.H{
			"type":         trashType,
			"trash_can_id": trashcanID,
			"status":       status,
			"message":      msg,
		}

		c.JSON(http.StatusInternalServerError, result)

	}

}
