package admin

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/jwt"
	"eatplek/pkg/services"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Id        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Password string `json:"-"`

	Id        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	Token     string    `json:"token"`
}

var collection = db.Client.Database("eatplek").Collection("admin")

func GetAdminByUsername(username string) (Admin, error) {
	var admin Admin
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&admin)
	if err != nil {
		return admin, errors.New("username not found")
	}
	return admin, nil
}

func Register(c *gin.Context) (Admin, error) {
	var u Admin
	if err := c.ShouldBindJSON(&u); err != nil {
		return u, err
	}

	_, err := GetAdminByUsername(u.Username)
	if err == nil {
		return u, errors.New("username already exists")
	}

	u.Password, _ = services.HashPassword(u.Password)
	u.Id = services.GenerateId()
	u.CreatedOn = time.Now()

	_, err = collection.InsertOne(context.TODO(), u)

	if err != nil {
		return u, err
	}

	return u, nil
}

func Login(c *gin.Context) (LoginResponse, error) {
	var a Admin
	var u LoginResponse

	if err := c.ShouldBindJSON(&a); err != nil {
		return u, err
	}

	err := collection.FindOne(context.TODO(), bson.M{"username": a.Username}).Decode(&u)
	if err != nil {
		return u, errors.New("username not found")
	}

	if !services.CheckPasswordHash(a.Password, u.Password) {
		return u, errors.New("wrong password")
	}

	u.Token, err = jwt.GenerateToken(u.Id, "admin")
	if err != nil {
		return u, errors.New("error generating token")
	}

	return u, nil
}
