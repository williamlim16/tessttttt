package main

import (
	_ "encoding/base64"
	_ "encoding/json"
	_ "strings"
	"trash-separator/config"
	"trash-separator/controllers"

	_ "encoding/base64"
	_ "encoding/json"
	_ "strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type InDB struct {
	DB *gorm.DB
}

func main() {
	db := config.DBInit()

	inDB := &controllers.InDB{DB: db}

	router := gin.Default()

	// Microcontroller -> Server
	router.POST("/node/sendLog", inDB.SendLog)
	router.POST("/node/sendCapacity", inDB.SendCapacity)

	// Server -> Web
	router.GET("/node/getCapacity/:trash_can_id", inDB.GetSingleTrashCanCapacity)

	router.GET("/node/getLogs/", inDB.GetAllTrashCanLogs)

	router.Run("localhost:8888")
}
