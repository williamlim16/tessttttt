package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"

	// "net/url"
	"strconv"
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
	resultQueryUser := idb.DB.Table("user").Select("*").Where("email = ?", email).First(&user)
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
		//put into redis, with expiry of 14 days. [key = token, value = user Struct]
		log.Printf("usertoken = %s", userToken)
		userJSON, err := json.Marshal(user)
		if err != nil {
			result = gin.H{"status": "error", "msg": fmt.Sprintf("Error parsing user into JSON, user: %v", user)}
			c.JSON(http.StatusInternalServerError, result)
			return
		}
		redisError := idb.RedisClient.Set(userToken, userJSON, time.Hour*24*14)
		if redisError.Err() != nil {
			result = gin.H{"status": "error", "msg": fmt.Sprintf("Internal Redis DB error, error: %v", redisError.Err().Error())}
			c.JSON(http.StatusInternalServerError, result)
			return
		}
		//put into user cookie then redirect
		c.SetCookie("user_token", userToken, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		c.SetCookie("user_id", strconv.Itoa(user.Id), 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		c.SetCookie("user_name", user.Name, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		c.SetCookie("user_address", user.Address, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		c.SetCookie("user_company", user.Company_name, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)
		c.SetCookie("user_telephone", user.Telephone, 3600*24*14, "/", c.Request.URL.Hostname(), false, true)

		//redirection handled by frontend
		// location := url.URL{Path: "/"}
		// c.Redirect(http.StatusFound, location.RequestURI())
		result = gin.H{
			"status": "success",
			"msg":    "login successful",
		}
		c.JSON(http.StatusOK, result)

	} else { //wrong password
		result = gin.H{"status": "not authorized", "msg": "Invalid email and/or password combination"}
		c.JSON(http.StatusUnauthorized, result)
		return
	}
}

func (idb *InDB) AuthRegister(c *gin.Context) {
	var (
		userInput   structs.User
		confirmPass string

		result gin.H
	)
	userInput = structs.User{
		Name:         c.PostForm("name"),
		Email:        c.PostForm("email"),
		Telephone:    c.PostForm("telephone"),
		Address:      c.PostForm("address"),
		Company_name: c.PostForm("company_name"),
		Password:     c.PostForm("password"),
	}
	confirmPass = c.PostForm("confirm_password")

	phoneRegex := "^[0-9]{8,13}$" //number only, with length of 8 - 13
	re := regexp.MustCompile(phoneRegex)
	_, errEmail := mail.ParseAddress(userInput.Email)
	//validate input
	if userInput.Name == "" || errEmail != nil || !re.MatchString(userInput.Telephone) || userInput.Address == "" || userInput.Company_name == "" || userInput.Password == "" || confirmPass == "" {
		result = gin.H{"status": "error", "msg": "Invalid inputs!"}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	//match password
	if userInput.Password != confirmPass {
		result = gin.H{"status": "error", "msg": "Password and confirm password input is not the same!"}
		c.JSON(http.StatusBadRequest, result)
		return
	}

	//push to db
	resultInsertUser := idb.DB.Table("users").Create(&userInput)

	if resultInsertUser.Error == nil {
		result = gin.H{"status": "success", "msg": "Register successful, please redirect user"}
		c.JSON(http.StatusOK, result)
		return
	} else {
		result = gin.H{"status": "error", "msg": fmt.Sprintf("internal db error in insert user, error: %v", resultInsertUser.Error.Error())}
		c.JSON(http.StatusInternalServerError, result)
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
