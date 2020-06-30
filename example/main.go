package main

import (
	"fmt"
	"time"

	// "github.com/davecgh/go-spew/spew"
	ginzap "github.com/bijenkins/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.New()

	//logger, _ := zap.NewProduction()
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{
		"/tmp/zaptest.log",
	}
	logger, _ := config.Build()

	defer logger.Sync()
	//spew.Dump(config)

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Example ping request.
	r.POST("/ping", func(c *gin.Context) {

		user := User{}
		err := c.ShouldBindJSON(&user)
		if err != nil {
			//logger.Info(err.Error())
			c.JSON(400, err.Error())
			//panic(err.Error())
			return
		}

		// c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
		c.JSON(200, gin.H{"data": user})
	})

	// Example when panic happen.
	r.GET("/panic", func(c *gin.Context) {
		panic("An unexpected error happen!")
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
