package restaurant

import (
	"eatplek/pkg/jwt"
	"eatplek/pkg/restaurant"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := restaurant.Add(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "restaurant added",
		"result":  result,
	})
}

func Login(c *gin.Context) {
	result, err := restaurant.Login(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "login successfull",
		"result":  result,
	})
}

func Get(c *gin.Context) {
	result, err := restaurant.Get()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":     "restaurants found",
		"restaurants": result,
	})
}

func GetAll(c *gin.Context) {
	result, err := restaurant.GetAll()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":     "restaurants found",
		"restaurants": result,
	})
}

func GetOne(c *gin.Context) {
	result, err := restaurant.GetOne(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message":    "restaurant found",
		"restaurant": result,
	})
}

func Update(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := restaurant.Update(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "restaurant updated",
	})
}

func Delete(c *gin.Context) {
	if !jwt.IsAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := restaurant.Delete(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "restaurant deleted",
	})
}

func UpdateOpenStatus(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := restaurant.UpdateOpen(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "restaurant open status updated",
	})
}

func GetOpenStatus(c *gin.Context) {
	open,err := restaurant.GetOpenStatus(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"open" : open,
	})
}

func ChangePassword(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	err := restaurant.ChangePassword(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "password changed",
	})
}

func ResetPassword(c *gin.Context) {
	err := restaurant.ResetPassword(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "password reset link sent",
	})
}

func ResetAndChangePassword(c *gin.Context) {
    if jwt.IsTokenExpired(c.Request.Header["Token"][0]) {
        c.JSON(400, gin.H{
            "message": "invalid token",
        })
        return
    }

	err := restaurant.ResetAndChangePassword(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "password reset successfully",
	})
}

func Profile(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	id := jwt.GetUserID(c.Request.Header["Token"][0])
	result, err := restaurant.Profile(id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "profile found",
		"profile": result,
	})
}

func UpdateProfile(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := restaurant.UpdateProfile(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "profile updated",
		"profile": result,
	})
}

func GetTimings(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	id := jwt.GetUserID(c.Request.Header["Token"][0])

	result, err := restaurant.GetTimings(id)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "timings found",
		"timings": result,
	})
}

func UpdateTimings(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{
			"message": "unauthorized",
		})

		return
	}

	result, err := restaurant.UpdateTimings(c)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "timings updated",
		"timings": result,
	})
}

func UpdateDeviceToken(c *gin.Context) {
	if !jwt.IsHotelAdmin(c.Request.Header["Token"][0]){
		c.JSON(401,gin.H{
			"message":"unauthorized",
		})

		return
	}

	err := restaurant.UpdateDeviceToken(c)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "token updated",
	})

}
