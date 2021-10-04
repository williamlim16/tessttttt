package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"trash-separator/structs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		fmt.Println("ERROR: ", resultInsertLog.Error)

		status = "error"
		msg = "Log insertion failed"

		result = gin.H{
			"status":  status,
			"message": msg,
		}

		c.JSON(http.StatusInternalServerError, result)

	}
}

func (idb *InDB) SendCapacity(c *gin.Context) {
	var (
		result gin.H
		status string
		msg    string
	)

	trashcanID, _ := strconv.Atoi(c.PostForm("trash_can_id"))
	organicCapacity, _ := strconv.Atoi(c.PostForm("organic_capacity"))
	anorganicCapacity, _ := strconv.Atoi(c.PostForm("anorganicCapacity"))

	insertCapacity := structs.Trash_capacity{
		Trash_id:           trashcanID,
		Organic_capacity:   organicCapacity,
		Anorganic_capacity: anorganicCapacity,
		Created_at:         time.Now(),
	}

	resultInsertCapacity := idb.DB.Table("trash_capacity").Create(&insertCapacity)

	if resultInsertCapacity.Error == nil {
		status = "success"
		msg = "Capacity log successfully added"

		result = gin.H{
			"status":  status,
			"message": msg,
		}

		c.JSON(http.StatusOK, result)
	} else {
		fmt.Println("ERROR: ", resultInsertCapacity.Error)

		status = "error"
		msg = "Capacity Log insertion failed"

		result = gin.H{
			"status":  status,
			"message": msg,
		}

		c.JSON(http.StatusInternalServerError, result)
	}
}

func (idb *InDB) GetSingleTrashCanCapacity(c *gin.Context) {
	var (
		trashCapacity structs.TrashCapacity
		status        string
		msg           string
		result        gin.H
	)

	trashCanID := c.Param("trash_can_id")

	resultGetSingleTrashCapacity := idb.DB.Table("trash_capacity").
		Select("trash_capacity.*, trash_version.organic_max_height, trash_version.inorganic_max_height").
		Joins("left join trash on trash.id = trash_capacity.trash_id").
		Joins("left join trash_version on trash_version.id = trash.trash_version_id").
		Where("trash_capacity.trash_id = ?", trashCanID).
		Find(&trashCapacity)

	if resultGetSingleTrashCapacity.Error == nil {

		status = "success"
		msg = "Get single trash can capacity successful"

		result = gin.H{
			"status":               status,
			"msg":                  msg,
			"trash_can_id":         trashCapacity.Trash_can_id,
			"organic_capacity":     trashCapacity.Organic_capacity,
			"anorganic_capacity":   trashCapacity.Anorganic_capacity,
			"organic_max_height":   trashCapacity.Organic_max_height,
			"anorganic_max_height": trashCapacity.Anorganic_max_height,
		}

		c.JSON(http.StatusOK, result)

	} else {
		status = "error"
		msg = resultGetSingleTrashCapacity.Error.Error()

		result = gin.H{
			"status": status,
			"msg":    msg,
		}

		c.JSON(http.StatusInternalServerError, result)

	}
}

func (idb *InDB) GetAllTrashCanLogs(c *gin.Context) {
	var (
		trashLogs     []structs.TrashLogs
		trashReadings map[int][]structs.Trash_reading
		status        string
		msg           string
		result        gin.H
	)
	trashReadings = map[int][]structs.Trash_reading{}
	rows, err := idb.DB.Table("trash_reading").Select("*").Rows()
	if idb.DB.Error == gorm.ErrRecordNotFound {
		result = gin.H{"status": "success", "msg": "no data"}
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	if err != nil {
		result = gin.H{"status": "error", "msg": "internal db error"}
		c.JSON(http.StatusInternalServerError, result)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var reading structs.Trash_reading
		idb.DB.ScanRows(rows, &reading)

		if idb.DB.Error != nil {
			result = gin.H{"status": "error", "msg": "internal db error"}
			c.JSON(http.StatusInternalServerError, result)
			return
		}
		trashReadings[reading.Trash_id] = append(trashReadings[reading.Trash_id], reading)
	}

	for key, val := range trashReadings {
		var temp structs.TrashLogs
		temp.Trash_can_id = key
		temp.Trash_reading = val
		trashLogs = append(trashLogs, temp)
	}
	status = "fetch log ok"
	msg = "ok"
	result = gin.H{
		"status": status,
		"msg":    msg,
		"data":   trashLogs,
	}

	c.JSON(http.StatusOK, result)

}
