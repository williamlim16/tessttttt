package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"trash-separator/structs"

	"github.com/gin-gonic/gin"
)

func (idb *InDB) SendLog(c *gin.Context) {
	var (
		result gin.H
		status string
		msg    string
	)

	trashcanID, _ := strconv.Atoi(c.PostForm("trash_can_id"))
	category := c.PostForm("category")
	trashType := c.PostForm("type")

	fmt.Println("Trash Type:", trashType)

	insertLog := structs.Trash_reading{
		Trash_id:   trashcanID,
		Category:   category,
		Type:       trashType,
		Created_at: time.Now(),
	}

	resultInsertLog := idb.DB.Table("trash_reading").Create(&insertLog)

	if resultInsertLog.Error == nil {
		status = "success"
		msg = "Log successfully added"

		result = gin.H{
			"status":  status,
			"message": msg,
		}

		c.JSON(http.StatusOK, result)

	} else {
		status = "error"
		msg = "Log insertion failed"

		result = gin.H{
			"status":  status,
			"message": msg,
		}

		c.JSON(http.StatusInternalServerError, result)

	}

}
