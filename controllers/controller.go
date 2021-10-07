package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"trash-separator/structs"
	"trash-separator/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (idb *InDB) SendLog(c *gin.Context) {
	var (
		result gin.H
		status string
		msg    string
	)

	trashcanID, err := strconv.Atoi(c.PostForm("trash_can_id"))
	if err != nil {
		result = gin.H{"status": "error", "msg": "invalid trash id format"}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	category := c.PostForm("category")
	if category != "inorganic" && category != "organic" {
		result = gin.H{"status": "error", "msg": "invalid trash category"}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	trashType := c.PostForm("type")

	fmt.Println("Trash Type:", trashType)

	insertLog := structs.Trash_reading{
		Trash_id:     trashcanID,
		Category:     category,
		Type:         trashType,
		Created_date: time.Now(),
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
	inorganicCapacity, _ := strconv.Atoi(c.PostForm("inorganic_capacity"))

	insertCapacity := structs.Trash_capacity{
		Trash_id:           trashcanID,
		Organic_capacity:   organicCapacity,
		Inorganic_capacity: inorganicCapacity,
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
		Last(&trashCapacity)

	if resultGetSingleTrashCapacity.Error == nil {
		status = "success"
		msg = "Get single trash can capacity successful"

		result = gin.H{
			"status": status,
			"msg":    msg,
			"data":   trashCapacity,
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
	if len(trashLogs) == 0 {
		msg = "empty logs"
	}
	result = gin.H{
		"status": status,
		"msg":    msg,
		"data":   trashLogs,
	}

	c.JSON(http.StatusOK, result)

}

func (idb *InDB) GetSingleTrashCanLogs(c *gin.Context) {
	var (
		trashLogs        structs.TrashLogs
		trashReadings    []structs.Trash_reading
		trashCanIDString string
		trashCanID       int
		status           string
		msg              string
		result           gin.H
		err              error
	)
	trashCanIDString = c.Param("trash_can_id")
	trashCanID, err = strconv.Atoi(trashCanIDString)
	if err != nil {
		result = gin.H{"status": "error", "msg": "invalid trash id format"}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	trashReadings = []structs.Trash_reading{}
	rows, err := idb.DB.Table("trash_reading").Select("*").Where("trash_id = ?", trashCanID).Rows()
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
		trashReadings = append(trashReadings, reading)
	}

	trashLogs.Trash_can_id = trashCanID
	trashLogs.Trash_reading = trashReadings

	status = "fetch log ok"
	msg = "ok"
	if len(trashLogs.Trash_reading) == 0 {
		msg = "empty logs"
	}
	result = gin.H{
		"status": status,
		"msg":    msg,
		"data":   trashLogs,
	}

	c.JSON(http.StatusOK, result)

}

func (idb *InDB) GetTopTrashCans(c *gin.Context) {
	var (
		trashReading []structs.TrashReading
		status       string
		msg          string
		userId       string
		result       gin.H
	)

	tn := time.Now()
	year, week := tn.ISOWeek()
	firstDayOfWeek := util.WeekStart(year, week)
	userId = getUserIdFromRedis(idb, c)

	log.Printf("User id : %s", userId)
	resultGetTopTrashCans := idb.DB.Table("trash_reading").
		Select("trash.trash_code AS Trash_sorter_name, trash.location AS Trash_sorter_location, count(*) as Total").
		Joins("left join trash on trash.id = trash_reading.trash_id").
		Group("trash.id").
		Where("trash_reading.created_date BETWEEN ? AND ?", firstDayOfWeek, tn).
		Where("trash.user_id = ?", userId).
		Order("Total desc").
		Limit(5).
		Find(&trashReading)

	if resultGetTopTrashCans.Error == nil {
		status = "success"
		msg = "Successfully get current week data"

		result = gin.H{
			"user_id": userId,
			"data":    trashReading,
			"status":  status,
			"msg":     msg,
		}

		c.JSON(http.StatusOK, result)

	} else {
		status = "error"
		msg = fmt.Sprintf("db error: %v", resultGetTopTrashCans.Error.Error())
		result = gin.H{
			"status": status,
			"msg":    msg,
		}

		c.JSON(http.StatusInternalServerError, result)
	}

}

func getUserIdFromRedis(idb *InDB, c *gin.Context) string {
	userToken, _ := c.Cookie("user_token")
	redisResp, err := idb.RedisClient.Get(userToken).Result()
	if err != nil {
		return ""
	}
	user := structs.User{}
	json.Unmarshal([]byte(redisResp), &user)
	return strconv.Itoa(user.Id)
}
