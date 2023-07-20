package main

import (
	"eatplek/pkg/db"
	"eatplek/pkg/handlers"
	"eatplek/pkg/handlers/admin"
	"eatplek/pkg/handlers/cart"
	"eatplek/pkg/handlers/category"
	"eatplek/pkg/handlers/feedback"
	"eatplek/pkg/handlers/food"
	"eatplek/pkg/handlers/orders"
	"eatplek/pkg/handlers/requests"
	"eatplek/pkg/handlers/restaurant"
	"eatplek/pkg/handlers/revenue"
	"eatplek/pkg/handlers/user"

	"github.com/gin-gonic/gin"
    "os"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main () {
	db.Connect()
	handler()
}

func handler() {
	r := gin.New()
	r.Use(CORSMiddleware())

	r.POST("/admin/register", admin.Register)
	r.POST("/admin/login", admin.Login)

	r.POST("/restaurant", restaurant.Add)
	r.POST("/restaurant/login", restaurant.Login)
	r.GET("/restaurant", restaurant.Get)
	r.GET("/restaurant/all", restaurant.GetAll)
	r.GET("/restaurant/profile", restaurant.Profile)
	r.PUT("/restaurant/profile", restaurant.UpdateProfile)
	r.GET("/restaurant/timings", restaurant.GetTimings)
	r.GET("/restaurant/requests", cart.Requests)
	r.PUT("/restaurant/timings", restaurant.UpdateTimings)
	r.GET("/restaurant/:id", restaurant.GetOne)
	r.PUT("/restaurant/:id", restaurant.Update)
	r.PUT("/restaurant/status/:id", restaurant.UpdateOpenStatus)
	r.GET("/restaurant/status/:id", restaurant.GetOpenStatus)
	r.PUT("/restaurant/password", restaurant.ChangePassword)
    r.POST("/restaurant/password/reset", restaurant.ResetPassword)
    r.PUT("/restaurant/password/reset", restaurant.ResetAndChangePassword)
	r.PATCH("/restaurant/token",restaurant.UpdateDeviceToken)
	r.DELETE("/restaurant/:id", restaurant.Delete)

	r.POST("/category", category.Add)
	r.GET("/category", category.Get)
	r.PUT("/category/:id", category.Update)
	r.DELETE("/category/:id", category.Delete)

	r.POST("/food", food.Add)
	r.PUT("/food/:id", food.Update)
	r.PUT("/food/availability/:id", food.UpdateAvailability)
	r.GET("/food", food.Get)
	r.GET("/food/:id", food.GetOne)
	r.GET("/food/filter/:filter/:id", food.GetByFilter)
	r.DELETE("/food/:id", food.Delete)
	r.GET("/restaurant/food", food.MenuChange)
	r.PUT("/restaurant/food", food.UpdateMenu)

	r.POST("/cart", cart.Initialize)
	r.PUT("/cart", cart.Add)
	r.PUT("/cart/status", cart.UpdateStatus)
	r.GET("/cart/:userId", cart.Get)

	r.POST("/order", orders.CreateOrder)
	r.GET("/order", orders.GetOrders)
	r.GET("/order/:id", orders.GetOrderById)
	r.GET("/order/filter/:type/:id", orders.GetOrder)
    r.PUT("/order/status", orders.UpdateStatus)

    r.POST("/user", user.SendOTP)
    //r.POST("/user/register", user.Register)
    //r.POST("/user/verify", user.Verify)
	r.POST("/user/login", user.Login)
	r.PUT("/user", user.UpdateUser)
	r.GET("/user/:id", user.GetUser)

	r.POST("/request", requests.New)
	r.GET("/request", requests.Get)
	r.GET("/request/:id", requests.GetOne)
	r.DELETE("/request/:id", requests.Delete)

	r.POST("/feedback", feedback.Create)
	r.GET("/feedback", feedback.Get)
	r.GET("/feedback/:id", feedback.GetOne)
	r.DELETE("/feedback/:id", feedback.Delete)

	r.POST("/admin/revenue", revenue.GetRevenue)
	r.POST("/hotel/revenue", revenue.GetHotelRevenue)

	r.POST("/upload", handlers.Upload)

    port := os.Getenv("PORT")
    if port == "" {
        port = "9000"
    }

    r.Run("127.0.0.1:" + port)
}
