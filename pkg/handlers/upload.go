package handlers

import (
	"github.com/gin-gonic/gin"
	"eatplek/pkg/services"
	"eatplek/pkg/jwt"
)

func Upload(c *gin.Context){
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result,err := services.UploadFile(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"link": result,"message":"file uploaded successfully"})
}
