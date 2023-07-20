package revenue

import (
	"eatplek/pkg/jwt"
	"eatplek/pkg/revenue"

	"github.com/gin-gonic/gin"
)

func GetRevenue(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	result, err := revenue.GetRevenue(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, result)
}

func GetHotelRevenue(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	result, err := revenue.GetHotelRevenue(c)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, result)
}
