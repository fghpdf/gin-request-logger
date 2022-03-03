# gin-request-logger
Request Logger Middleware for Gin Framework

This is a middleware for [Gin](https://github.com/gin-gonic/gin) framework.

It uses [zap](https://github.com/uber-go/zap) to provide a log featurues in package.

Fork from https://github.com/karmadon/gin-request-logger

Like [This project](https://github.com/gin-contrib/zap)

## Usage


### Install
Download and install using [go module](https://blog.golang.org/using-go-modules):
```shell
go get github.com/fghpdf/gin-request-logger
```

### Usage
```golang
package main

import (
  "fmt"
  "time"

  gin_request_logger "github.com/fghpdf/gin-request-logger"
  "github.com/gin-gonic/gin"
  "go.uber.org/zap"
)

func main() {
  r := gin.New()

  logger, _ := zap.NewProduction()

  r.Use(gin_request_logger.New(gin_request_logger.Options{
    LogResponse: true,
    Logger: logger
  }))

  // Example ping request.
  r.GET("/ping", func(c *gin.Context) {
    c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Example when panic happen.
  r.GET("/panic", func(c *gin.Context) {
    panic("An unexpected error happen!")
  })
  
  // Example post request
  r.GET("/post", func(c *gin.Context) {
    c.String(200, "Hello!")
  })

  // Listen and Server in 0.0.0.0:8080
  r.Run(":8080")
}
```
