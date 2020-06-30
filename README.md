# zap

[![Build Status](https://travis-ci.org/gin-contrib/zap.svg?branch=master)](https://travis-ci.org/gin-contrib/zap) [![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/zap)](https://goreportcard.com/report/github.com/gin-contrib/zap)
[![GoDoc](https://godoc.org/github.com/gin-contrib/zap?status.svg)](https://godoc.org/github.com/gin-contrib/zap)
[![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin)

Alternative logging through [zap](https://github.com/uber-go/zap). Thanks for [Pull Request](https://github.com/gin-gonic/contrib/pull/129) from [@yezooz](https://github.com/yezooz)

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/bijenkins/zap
```

Import it in your code:

```go
import "github.com/bijenkins/zap"
```

## Example

See the [example](example/main.go).

[embedmd]:# (example/main.go go)
```go
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
	Username string `json:"username"   binding:"required"`
	Password string `json:"password"   binding:"required"`
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
```
