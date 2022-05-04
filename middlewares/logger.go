package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

func Logger(v ...url.Values) gin.HandlerFunc {
	return gin.Logger()
}

func Recovery(v ...url.Values) gin.HandlerFunc {
	return gin.Recovery()
}