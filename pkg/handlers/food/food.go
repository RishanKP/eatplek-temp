package food

import (
	"eatplek/pkg/food"
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

	result, err := food.Add(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food added",
		"result":  result,
	})
}

func getByFilterRestaurant(c *gin.Context) {
	id := c.Param("id")
	category := c.Query("category")

    
	if c.Query("foredit") == "true" {
		result, err := food.GetByRestaurantForEdit(id)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}

		c.JSON(200, gin.H{
			"message": "foods found",
			"foods":   result,
		})

		return
	}

    if category == "" {
        if c.Query("usertype") == "admin" {
            result, err := food.GetByRestaurant(id,"admin")
            if err != nil {
                c.JSON(400, gin.H{
                    "message": err.Error(),
                })

                return
            }

            c.JSON(200, gin.H{
                "message": "foods found",
                "foods":   result,
            })

            return
        }

		result, err := food.GetByRestaurant(id,"")
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})

			return
		}

		c.JSON(200, gin.H{
			"message": "foods found",
			"foods":   result,
		})

		return
	}

	result, err := food.GetByRestaurantAndCategory(id, category)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "foods found",
		"foods":   result,
	})
}

func getByFilterCategory(c *gin.Context) {
	category := c.Param("id")

	result, err := food.GetByCategory(category)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "foods found",
		"foods":   result,
	})
}

func GetByFilter(c *gin.Context) {
	filter := c.Param("filter")
	switch filter {
	case "restaurant":
		getByFilterRestaurant(c)
	case "category":
		getByFilterCategory(c)
	}
}

func Get(c *gin.Context) {
	result, err := food.Get()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "foods found",
		"foods":   result,
	})
}

func GetOne(c *gin.Context) {
	result, err := food.GetOne(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food found",
		"food":    result,
	})
}

func Update(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := food.Update(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food updated",
	})
}

func Delete(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := food.Delete(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food deleted",
	})
}

func UpdateAvailability(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := food.UpdateAvailability(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food availability updated",
	})
}

func MenuChange(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}
	id := jwt.GetUserID(c.Request.Header["Token"][0])

	result, err := food.MenuChange(id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "foods found",
		"foods":   result,
	})
}

func UpdateMenu(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	id := jwt.GetUserID(c.Request.Header["Token"][0])

	update, err := food.UpdateMenu(c, id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "food menu updated",
		"update":  update,
	})
}
