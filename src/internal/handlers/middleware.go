package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	userContext = "userID"
)

func (h *Handlers) userIndentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)

	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
	}

	c.Set(userContext, id)
}

func (h *Handlers) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")

	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getUserId(c *gin.Context) (primitive.ObjectID, error) {
	return getIdByContext(c, userContext)
}

func getIdByContext(c *gin.Context, context string) (primitive.ObjectID, error) {
	idFromContext, ok := c.Get(context)
	if !ok {
		return primitive.ObjectID{}, errors.New("user context not found")
	}

	id, err := primitive.ObjectIDFromHex(idFromContext.(string))
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (h *Handlers) CorsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
