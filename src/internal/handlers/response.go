package handlers

import (
	"github.com/avalonprod/gasstrem/src/packages/logger"
	"github.com/gin-gonic/gin"
)

type responseData struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

type responseId struct {
	ID interface{} `json:"id"`
}

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	logger.Error(message)
	c.AbortWithStatusJSON(statusCode, response{message})
}
