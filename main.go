package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"os"
)

var err error

const recordBasePath = "/var/log/record/"

func setupRouter() *gin.Engine {

	gin.DisableConsoleColor()
	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"root": "root",
		"admin": "admin",
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authorized.GET("/play/:id", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		callId := c.Params.ByName("id")

		wavMediaFile := recordBasePath + callId + ".wav"
		_, err := os.Stat(wavMediaFile)

		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"user": user, "message": "Not Found", "callId": callId})
		} else  if err != nil{
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Error", "reason" :  err.Error()})
		}  else {
			c.File(wavMediaFile)
		}
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":9090")
}
