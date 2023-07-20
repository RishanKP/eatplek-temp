package category

import (
	"eatplek/pkg/category"
	"eatplek/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := category.Add(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "category added",
		"result":  result,
	})
}

func Get(c *gin.Context) {
	result, err := category.Get()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":    "categories found",
		"categories": result,
	})
}

func Update(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := category.Update(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "category updated",
	})
}

func Delete(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := category.Delete(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "category deleted",
	})
}
