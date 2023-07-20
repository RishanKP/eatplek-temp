package requests

import (
	"eatplek/pkg/jwt"
	"eatplek/pkg/requests"

	"github.com/gin-gonic/gin"
)

func New(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := requests.New(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result":  result,
		"message": "request created",
	})
}

func Get(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := requests.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result":  result,
		"message": "requests retrieved",
	})
}

func GetOne(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := requests.GetOne(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"result":  result,
		"message": "request retrieved",
	})
}

func Delete(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := requests.Delete(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "request deleted",
	})
}
