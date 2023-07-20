package cart

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/user"
	"eatplek/pkg/restaurant"
	"eatplek/pkg/notification"
	"eatplek/pkg/services"
	"errors"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Item struct {
	FoodId   string  `json:"food_id"`
	Name     string  `json:"name"`
	Image    string  `json:"image"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`

	Total float64 `json:"total"`
}

type CartItem struct {
	UserId string `json:"user_id"`
	Item   Item   `json:"item"`
}

type Cart struct {
	UserId         string  `json:"user_id"`
	RestaurantId   string  `json:"restaurant_id"`
	RestaurantName string  `json:"restaurant_name"`
	Items          []Item  `json:"items"`
	NumberOfGuests int     `json:"number_of_guests"`
	DeviceToken    string  `json:"device_token"`
	Status         int     `json:"status"` //-1: pending, 0: approved, 1: declined
	Time           string  `json:"time"`
	Type           string  `json:"type"` // dine in , take away
	TotalAmount    float64 `json:"total_amount"`

	UpdatedOn time.Time `json:"updated_on"`
}

type Request struct {
	NumberOfGuests int     `json:"number_of_guests"`
	Time           string  `json:"time"`
	Type           string  `json:"type"` // dine in , take away
	UserId         string  `json:"user_id"`
}

type DeclinedOrder struct {
	Restaurant string    `json:"restaurant"`
	Customer   string    `json:"customer"`
	Contact    string    `json:"contact"`
	DeclinedAt time.Time `json:"declined_at"`
	Type       string    `json:"type"`
}

var collection = db.Client.Database("eatplek").Collection("cart")

func Initialize(c *gin.Context) (Cart, error) {
	var cart Cart
	c.BindJSON(&cart)
    
    openingtime,closingtime,err := restaurant.GetOpeningAndClosingTime(cart.RestaurantId)
    if err != nil{
        return cart,err
    }

    err = services.ValidateRequestTime(openingtime,closingtime,cart.Time)
    if err != nil{
        notification.NotifyUser("Invalid Time",err.Error(),cart.DeviceToken)
        return cart,err
    }

	cart.UpdatedOn = time.Now()
	cart.Status = -1
	cart.Items = []Item{}

	_, err = collection.DeleteOne(context.TODO(), bson.D{
		{"userid", cart.UserId},
	})
	if err != nil {
		return cart, err
	}
	_, err = collection.InsertOne(context.TODO(), cart)
	if err != nil {
		return cart, err
	}

	restaurant_tokens,restaurant_web_tokens, err := restaurant.GetDeviceToken(cart.RestaurantId)
	if err != nil{
		fmt.Println("error fetching device token")
	}
	
	// send notification to restaurant to approve/decline order
	notification.NotifyRestaurant(cart.Type,cart.Time,cart.NumberOfGuests,restaurant_tokens,restaurant_web_tokens,cart.UserId)
	return cart, nil
}

func Add(c *gin.Context) (Item, error) {
	var item CartItem

	if err := c.ShouldBindJSON(&item); err != nil {
		return item.Item, err
	}

	item.Item.Total = item.Item.Price * float64(item.Item.Quantity)

	filter := bson.D{{"userid", item.UserId}}

	var cart Cart
	err := collection.FindOne(context.TODO(), filter).Decode(&cart)
	if err != nil {
		return item.Item, err
	}

	found := false
	for i, v := range cart.Items {
		if v.FoodId == item.Item.FoodId {
			found = true

			if v.Quantity+item.Item.Quantity == 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
				cart.TotalAmount -= v.Total
			} else {
				cart.Items[i].Quantity += item.Item.Quantity
				cart.Items[i].Total += item.Item.Total
				cart.TotalAmount += item.Item.Total
			}
			break
		}
	}

	if found {
		_, err = collection.UpdateOne(context.TODO(), filter, bson.D{
			{"$set", bson.D{
				{"items", cart.Items},
				{"totalamount", cart.TotalAmount},
				{"updatedon", time.Now()},
			}},
		})

		if err != nil {
			return item.Item, err
		}

		return item.Item, nil
	}

	if item.Item.Quantity <= 0 {
		return item.Item, errors.New("quantity must be greater than 0")
	}

	opts := options.Update().SetUpsert(true)
	update := bson.M{
		"$push": bson.M{"items": item.Item},
		"$set":  bson.M{"updatedon": time.Now()},
		"$inc":  bson.M{"totalamount": item.Item.Total},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return item.Item, err
	}

	return item.Item, nil
}

func GetCartByUserId(userid string) (Cart, error) {
	var cart Cart

	filter := bson.D{{"userid", userid}}
	err := collection.FindOne(context.TODO(), filter).Decode(&cart)
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func Get(userid string) (Cart, error) {
	var cart Cart
	filter := bson.D{{"userid", userid}}
	err := collection.FindOne(context.TODO(), filter).Decode(&cart)
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func UpdateStatus(c *gin.Context) error{
	type status struct{
		Userid string `json:"userid"`
		Status int `json:"status"` //0 approved, 1 rejected
	}

	var s status

	c.BindJSON(&s)
	_,err := collection.UpdateOne(context.TODO(),bson.D{{"userid",s.Userid}},bson.M{"$set":bson.D{{"status",s.Status},}})

	if err != nil{
		return err
	}
	
	var title,body string

	type token struct{
		DeviceToken string `bson:"devicetoken"`
	}
	var t token

	opts := options.FindOne().SetProjection(bson.D{{"devicetoken",1},})
	err = collection.FindOne(context.TODO(),bson.D{{"userid",s.Userid},},opts).Decode(&t)
	if err != nil{
		fmt.Println("error fetching user device token",err)
	}

	if s.Status == 1{
		title = "Request Rejected"
		body = "Oops..It looks like there are no seats available at the moment"
	}

	if s.Status == 0{
		title = "Request Approved"
		body = "Thank you for choosing us. Can't wait to serve you !!"
	}
    
    go func(){    
        if s.Status == 1 {
            var c Cart 
            err = collection.FindOne(context.TODO(),bson.D{{"userid",s.Userid}}).Decode(&c)
            if err != nil{
                fmt.Println("error finding document")
            }
            
            var d DeclinedOrder

            d.Restaurant = c.RestaurantName

            var u user.User
            u,err = user.GetUser(s.Userid)
            if err != nil{
                fmt.Println("error fetching user details")
            }

            d.Customer = u.Name
            d.Contact = u.Phone
            d.DeclinedAt = time.Now()
            d.Type = c.Type

	        d_collection := db.Client.Database("eatplek").Collection("declined_orders")
            _,err = d_collection.InsertOne(context.TODO(),d)
            if err != nil{
                fmt.Println("failed to insert document to declined orders")
            }
        }
    }()


	go notification.NotifyUser(title,body,t.DeviceToken)

	return nil
}

func Requests(rid string) ([]Request,error){
	var requests []Request
	curr,err := collection.Find(context.TODO(),bson.D{{"restaurantid",rid},{"status",-1}})

	if err != nil{
		return requests,err
	}

	for curr.Next(context.TODO()){
		var r Request

		err = curr.Decode(&r)
		if err != nil{
			return requests,err
		}

		requests = append(requests,r)
	}

	return requests,nil
}
