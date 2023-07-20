package orders

import (
	"context"
	"eatplek/pkg/cart"
	"eatplek/pkg/db"
	"eatplek/pkg/restaurant"
	"eatplek/pkg/services"
	"eatplek/pkg/user"
	"time"
    "fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection = db.Client.Database("eatplek").Collection("orders")

type Order struct {
	Id         string                        `json:"id"`
	Cart       cart.Cart                     `json:"cart"`
	User       user.UserForOrder             `json:"user"`
	Restaurant restaurant.RestaurantForOrder `json:"restaurant"`
	Comments   string                        `json:"comments"`
	Status     int                           `json:"status"`
	CreatedAt  time.Time                     `json:"created_at"`
}

type NewOrder struct {
	UserId   string `json:"user_id"`
	Comments string `json:"comments"`
}

type DeclinedOrder struct {
	Restaurant string    `json:"restaurant"`
	Customer   string    `json:"customer"`
	Contact    string    `json:"contact"`
	DeclinedAt time.Time `json:"declined_at"`
	Type       string    `json:"type"`
}

func CreateOrder(c *gin.Context) (Order, error) {
	var neworder NewOrder
	c.BindJSON(&neworder)

	var order Order

	var err error
	order.User, err = user.GetDetails(neworder.UserId)

	if err != nil {
		return Order{}, err
	}

	order.Cart, err = cart.GetCartByUserId(neworder.UserId)
	if err != nil {
		return Order{}, err
	}

	order.Restaurant, err = restaurant.GetDetails(order.Cart.RestaurantId)
	if err != nil {
		return Order{}, err
	}

	order.Comments = neworder.Comments
	order.Id = services.GenerateId()
	order.Status = -1
	order.CreatedAt = time.Now()
	_, err = collection.InsertOne(context.TODO(), order)
	if err != nil {
		return Order{}, err
	}

    go func(){
    msg := fmt.Sprintf("Your order has been placed successfully.\nTime: %s\nOrder Type: %s\nThank you for choosing eatplek.",order.Cart.Time,order.Cart.Type) 
        err = services.SendOTP(order.User.Phone,msg,"ORDER_TID")
        if err != nil{
            //do nothing
        }
    }()
    return order, nil
}

func GetOrders() ([]Order, error) {
	var orders []Order

	findOptions := options.Find().SetSort(bson.D{{"createdat", -1}})
	cur, err := collection.Find(context.TODO(), bson.D{}, findOptions)

	for cur.Next(context.TODO()) {
		var order Order
		err := cur.Decode(&order)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func GetOrdersByUserId(userid string) ([]Order, error) {
	var orders []Order

	opts := options.Find()
	opts.SetSort(bson.D{{"createdat", -1}})
	opts.SetProjection(bson.D{{"user", 0}})
	filter := bson.D{{"cart.userid", userid}}
	cur, err := collection.Find(context.TODO(), filter, opts)

	for cur.Next(context.TODO()) {
		var order Order
		err := cur.Decode(&order)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func GetOrdersByRestaurantId(rid string) ([]Order, error) {
	var orders []Order

	findOptions := options.Find().SetSort(bson.D{{"createdat", -1}})

	filter := bson.D{{"cart.restaurantid", rid}}
	cur, err := collection.Find(context.TODO(), filter, findOptions)

	for cur.Next(context.TODO()) {
		var order Order
		err := cur.Decode(&order)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func GetOrderById(id string) (Order, error) {
	var order Order
	filter := bson.D{{"id", id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&order)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func GetDeclinedOrders() ([]DeclinedOrder, error) {
	var orders []DeclinedOrder

	d_collection := db.Client.Database("eatplek").Collection("declined_orders")

	findOptions := options.Find().SetSort(bson.D{{"declineddat", -1}})
	cur, err := d_collection.Find(context.TODO(), bson.D{}, findOptions)

	for cur.Next(context.TODO()) {
		var order DeclinedOrder
		err := cur.Decode(&order)
		if err != nil {
			return []DeclinedOrder{}, err
		}
		orders = append(orders, order)
	}

	if err != nil {
		return []DeclinedOrder{}, err
	}

	return orders, nil
}

func UpdateStatus(c *gin.Context) error {
	type status struct {
		Status int    `json:"status"`
		Id     string `json:"id"`
	}
	var s status
	c.BindJSON(&s)

	filter := bson.D{{"id", s.Id}}
	update := bson.D{{"$set", bson.D{{"status", s.Status}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
