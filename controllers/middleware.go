package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"trash-separator/structs"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func (idb *InDB) EnableMiddleware() {
	idb.Middleware = true
}

func (idb *InDB) MWCheckUserTokenCookie() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		if idb.Middleware {
			log.Println("Executing cookie middleware")
			//check for cookie
			user := structs.User{}
			token, err := c.Cookie("user_token")
			if err != nil {
				result := gin.H{"status": "not authorized", "msg": "You are not logged in"}
				c.JSON(http.StatusUnauthorized, result)
				c.Abort()
				return
			}
			log.Printf("token = %s", token)

			//check if cookie is valid or not
			resJSON, err := idb.RedisClient.Get(token).Result()
			if err == redis.Nil {
				//not found, delete all cookie
				c.SetCookie("user_token", "", 0, "/", c.Request.URL.Hostname(), false, true)
				c.SetCookie("user_id", "", 0, "/", c.Request.URL.Hostname(), false, true)
				c.SetCookie("user_name", "", 0, "/", c.Request.URL.Hostname(), false, true)
				c.SetCookie("user_address", "", 0, "/", c.Request.URL.Hostname(), false, true)
				c.SetCookie("user_company", "", 0, "/", c.Request.URL.Hostname(), false, true)
				c.SetCookie("user_telephone", "", 0, "/", c.Request.URL.Hostname(), false, true)

				result := gin.H{"status": "not authorized", "msg": "You are not logged in"}
				c.JSON(http.StatusUnauthorized, result)
				c.Abort()
				return
			}
			json.Unmarshal([]byte(resJSON), &user)
			log.Printf("resJSON = %s", resJSON)

		}

		c.Next()

	})
}
