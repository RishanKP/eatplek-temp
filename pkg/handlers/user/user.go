package user

import (
	"eatplek/pkg/jwt"
	"eatplek/pkg/user"

	"github.com/gin-gonic/gin"
)

func SendOTP(c *gin.Context){
	err := user.SendOTP(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "OTP sent successfully"})
}
//func Register(c *gin.Context) {
//	err := user.Register(c)
//	if err != nil {
//		c.JSON(400, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(200, gin.H{"message": "User registered successfully, verification email sent"})
//}

func Verify(c *gin.Context) {
	err := user.Verify(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User verified successfully"})
}

func Login(c *gin.Context) {
	user, err := user.Login(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User logged in successfully", "user": user})
}

func Update(c *gin.Context) {
	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := user.Update(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func UpdateUser(c *gin.Context) {
	if jwt.IsHotelAdmin(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	if !jwt.IsValidToken(c.Request.Header["Token"][0]) {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	err := user.UpdateUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func GetUser(c *gin.Context) {
	user, err := user.GetUser(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User retrieved successfully", "user": user})
}
