package orders

import (
	"eatplek/pkg/jwt"
	"eatplek/pkg/orders"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := orders.CreateOrder(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": result, "message": "Order created successfully"})

	return
}

func GetOrders(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	status := c.Query("status")

	if status == "declined" {
		result, err := orders.GetDeclinedOrders()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": result, "message": "Orders fetched successfully"})

		return
	}

	result, err := orders.GetOrders()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": result, "message": "Orders fetched successfully"})

	return
}

func GetOrder(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	if c.Param("type") == "user" {
		result, err := orders.GetOrdersByUserId(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"result": result, "message": "Orders fetched successfully"})

		return
	}

	result, err := orders.GetOrdersByRestaurantId(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": result, "message": "Orders fetched successfully"})

	return
}

func GetOrderById(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := orders.GetOrderById(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"result": result, "message": "Order fetched successfully"})

	return
}

func UpdateStatus(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := orders.UpdateStatus(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Order status updated successfully"})

	return
}
