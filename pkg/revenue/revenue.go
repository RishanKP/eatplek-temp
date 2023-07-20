package revenue

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/jwt"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurants struct {
	Name    string
	Revenue float64
}

type Revenue struct {
	Total       float64
	Restaurants []Restaurants
}

type HotelRevenue struct {
	Total  float64
	Orders []Order
}

type Order struct {
	ID            string             `json:"id"            bson:"id"`
	CreatedAt     primitive.DateTime `json:"createdat"     bson:"createdat"`
	Name          string             `json:"name"          bson:"name"`
	Phone         string             `json:"phone"         bson:"phone"`
	TotalAmount   float64            `json:"totalamount"   bson:"totalamount"`
	NumberOfItems int32              `json:"numberofitems" bson:"numberofitems"`
}

type Dates struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func GetRevenue(c *gin.Context) (Revenue, error) {
	var d Dates
	c.BindJSON(&d)

	collection := db.Client.Database("eatplek").Collection("orders")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdat": bson.M{
					"$gte": d.StartDate,
					"$lte": d.EndDate,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$cart.restaurantid",
				"total": bson.M{
					"$sum": "$cart.totalamount",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "restaurant",
				"localField":   "_id",
				"foreignField": "id",
				"as":           "restaurant",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$restaurant",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"name":     "$restaurant.name",
				"customer": "$user.name",
				"contact":  "$user.contact",

				"revenue": "$total",
			},
		},
	}

	result, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return Revenue{}, err
	}

	var response []bson.M
	if err = result.All(context.TODO(), &response); err != nil {
		return Revenue{}, err
	}

	var revenue Revenue
	for _, v := range response {
		revenue.Total += v["revenue"].(float64)
		revenue.Restaurants = append(revenue.Restaurants, Restaurants{
			Name:    v["name"].(string),
			Revenue: v["revenue"].(float64),
		})
	}

	return revenue, nil
}

func GetHotelRevenue(c *gin.Context) (HotelRevenue, error) {
	type order struct {
		ID        string    `json:"id"        bson:"id"`
		CreatedAt time.Time `json:"createdat" bson:"createdat"`
		User      struct {
			Name  string `json:"name" bson:"name"`
			Phone string `json:"phone" bson:"phone"`
		}
		Cart struct {
			TotalAmount float64 `json:"totalamount" bson:"totalamount"`
		}
	}

	var d Dates
	c.BindJSON(&d)

	collection := db.Client.Database("eatplek").Collection("orders")

	id := jwt.GetUserID(c.Request.Header["Token"][0])

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"cart.restaurantid": id,
				"createdat": bson.M{
					"$gte": d.StartDate,
					"$lte": d.EndDate,
				},
			},
		},
		{
			"$sort": bson.M{
				"createdat": -1,
			},
		},
		{
			"$project": bson.M{
				"id":        "$id",
				"createdat": "$createdat",
				"user":      "$user.name",
				"phone":     "$user.phone",
				"total":     "$cart.totalamount",
				"numberoforders": bson.M{
					"$size": "$cart.items",
				},
			},
		},
	}

	result, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return HotelRevenue{}, err
	}

	var response []bson.M
	if err = result.All(context.TODO(), &response); err != nil {
		return HotelRevenue{}, err
	}

	fmt.Println(response)

	var revenue HotelRevenue

	for _, v := range response {
		revenue.Total += v["total"].(float64)
		revenue.Orders = append(revenue.Orders, Order{
			ID:            v["id"].(string),
			CreatedAt:     v["createdat"].(primitive.DateTime),
			Name:          v["user"].(string),
			Phone:         v["phone"].(string),
			TotalAmount:   v["total"].(float64),
			NumberOfItems: v["numberoforders"].(int32),
		})
	}

	return revenue, nil
}
