package feedback

import (
	"eatplek/pkg/feedback"
	"eatplek/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context) {
	err := feedback.Create(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Feedback created"})
}

func Get(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	feedbacks, err := feedback.Get()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, feedbacks)
}

func GetOne(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	feedback, err := feedback.GetOne(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, feedback)
}

func Delete(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	err := feedback.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Feedback deleted"})
}
