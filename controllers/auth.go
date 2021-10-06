package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"trash-separator/structs"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (idb *InDB) AuthLogin(c *gin.Context) {
	var (
		user      structs.User
		email     string
		password  string
		userToken string
		result    gin.H
	)

	email = c.PostForm("email")
	password = c.PostForm("password")

	if email == "" || password == "" {
		result = gin.H{"status": "error", "msg": "Failed parsing email and/or password"}
		c.JSON(http.StatusBadRequest, result)
		return
	}
	resultQueryUser := idb.DB.Table("user").Select("name, email, password").Where("email = ?", email).First(&user)
	if resultQueryUser.Error != nil { //some error
		if resultQueryUser.Error == gorm.ErrRecordNotFound { //got no record from db
			result = gin.H{"status": "not authorized", "msg": "Invalid email and/or password combination"}
			c.JSON(http.StatusUnauthorized, result)
			return
		} else { //other db errors
			result = gin.H{"status": "error", "msg": fmt.Sprintf("Internal DB error, error: %v", resultQueryUser.Error.Error())}
			c.JSON(http.StatusInternalServerError, result)
			return
		}
	}

	if user.Password == password { //same password
		//create random token
		userToken = randToken()
		//put into redis, with expiry of 14 days. [key = token, value = 1(valid), 0(invalid)]
		log.Printf("usertoken = %s", userToken)
		redisError := idb.RedisClient.Set(userToken, 1, time.Hour*24*14)
		if redisError.Err() != nil {
			result = gin.H{"status": "error", "msg": fmt.Sprintf("Internal Redis DB error, error: %v", redisError.Err().Error())}
			c.JSON(http.StatusInternalServerError, result)
			return
		}
		//put into user cookie then redirect
		c.SetCookie("user_token", userToken, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		location := url.URL{Path: "/"}
		c.Redirect(http.StatusFound, location.RequestURI())

	} else { //wrong password
		result = gin.H{"status": "not authorized", "msg": "Invalid email and/or password combination"}
		c.JSON(http.StatusUnauthorized, result)
		return
	}
}

func randToken() string {
	b := make([]byte, 16) //length is 16*2 characters
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (idb *InDB) NotImplemented(c *gin.Context) {
	res := gin.H{"status": "not implemented", "msg": "not implemented"}
	c.JSON(http.StatusOK, res)
}
