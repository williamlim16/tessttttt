package main

import (
	"trash-separator/config"
	"trash-separator/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InDB struct {
	DB *gorm.DB
}

func main() {
	db := config.DBInit()
	client := config.RedisInit()

	inDB := &controllers.InDB{DB: db, RedisClient: client}
	// inDB.EnableMiddleware()

	router := gin.Default()

	// Microcontroller -> Server
	router.POST("/node/sendLog", inDB.SendLog)
	router.POST("/node/sendCapacity", inDB.SendCapacity)

	// Server -> Web
	router.GET("/node/getCapacity/:trash_can_id", inDB.MWCheckUserTokenCookie(), inDB.GetSingleTrashCanCapacity)
	router.GET("/node/getTopTrashCans", inDB.MWCheckUserTokenCookie(), inDB.GetTopTrashCans)

	router.GET("/node/getLogs/", inDB.MWCheckUserTokenCookie(), inDB.GetAllTrashCanLogs)
	router.GET("/node/getLogs/:trash_can_id", inDB.MWCheckUserTokenCookie(), inDB.GetSingleTrashCanLogs)

	// Authentication
	router.POST("/api/login", inDB.AuthLogin)
	router.POST("/api/register", inDB.AuthRegister)
	router.GET("/", inDB.NotImplemented)

	router.Run("localhost:8888")
}
