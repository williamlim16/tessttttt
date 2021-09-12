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

	router.GET("/getPlasticSpaceCount", inDB.GetPlasticSpaceCount)
	router.GET("/getMetalSpaceCount", inDB.GetMetalSpaceCount)
	router.GET("/getGlassSpaceCount", inDB.GetGlassSpaceCount)
	router.POST("/sendLog", inDB.SendLog)

	router.Run("localhost:8888")
}
