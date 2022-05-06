package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zehongyang/ygins"
	"net/url"
)

func main() {
	ygins.Register(Login)
	ygins.Run()
}

func Login(v ...url.Values) gin.HandlerFunc {
	fmt.Println(v)
	return func(c *gin.Context) {
		
	}
}
