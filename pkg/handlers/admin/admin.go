package admin

import (
	"eatplek/pkg/admin"
	"eatplek/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	admin, err := admin.Register(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "admin created successfully", "user": admin})
}

func Login(c *gin.Context) {
	admin, err := admin.Login(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "admin logged in successfully", "user": admin})
}
