package requests

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Request struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	Restaurant   string  `json:"restaurant"`
	RestaurantId string  `json:"restaurant_id"`

	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

var collection = db.Client.Database("eatplek").Collection("requests")

func New(c *gin.Context) (Request, error) {
	var request Request
	c.BindJSON(&request)

	request.Id = services.GenerateId()
	request.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.Background(), request)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}

func Get() ([]Request, error) {
	var requests []Request

	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var request Request
		err := cur.Decode(&request)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func GetOne(id string) (Request, error) {
	var request Request
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&request)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}

func Delete(id string) error {
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}

	return nil
}
