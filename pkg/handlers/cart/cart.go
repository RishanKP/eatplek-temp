package cart

import (
	"eatplek/pkg/cart"
	"eatplek/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := cart.Add(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "item added to cart",
		"item":    result,
	})
}

func Get(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := cart.Get(c.Param("userId"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "cart retrieved",
		"cart":    result,
	})
}

func Initialize(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	result, err := cart.Initialize(c)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "cart initialized",
		"cart":    result,
	})
}

func UpdateStatus(c *gin.Context){
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]){
		c.JSON(401, gin.H{"error":"unauthorized"})
		return
	}

	err := cart.UpdateStatus(c)
	if err != nil{
		c.JSON(400, gin.H{
			"error":err,
		})
		return
	}

	c.JSON(200,gin.H{
		"message":"status updated",
	})
}

func Requests(c *gin.Context){
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]){
		c.JSON(401, gin.H{"error":"unauthorized"})
		return
	}
	
	id := jwt.GetUserID(c.Request.Header["Token"][0])
	result , err := cart.Requests(id)

	if err != nil{
		c.JSON(400, gin.H{
			"error":err,
		})
		return
	}

	c.JSON(200,gin.H{
		"requests": result,
		"message": "success",
	})
}
