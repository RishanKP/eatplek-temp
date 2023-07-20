package restaurant

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/food"
	"eatplek/pkg/jwt"
	"eatplek/pkg/services"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Restaurant struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsVeg    bool   `json:"isveg"`
	Capacity int    `json:"maximum_no_of_guests"`
	Open     bool   `json:"open"`
    	AC bool `json:"ac"`
    	NonAC bool `json:"non_ac"`

	DaysOpen  []string `json:"days_open"`
	OpenTime  string   `json:"opening_time"`
	CloseTime string   `json:"closing_time"`
	Type      string   `json:"type"`
	Image     string   `json:"image"`
	DineIn    bool     `json:"dine_in"`
	TakeAway  bool     `json:"take_away"`

	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

type HotelProfile struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type HotelTimings struct {
	DaysOpen  []string `json:"days_open"`
	OpenTime  string   `json:"opening_time"`
	CloseTime string   `json:"closing_time"`
	Open     bool   `json:"open"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RestaurantForOrder struct {
	Location string `json:"location"`
	Phone    string `json:"phone"`
}

type LoginResponse struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsVeg    bool   `json:"isveg"`
	Capacity int    `json:"maximum_no_of_guests"`

    Type     string `json:"type"`
    Image    string `json:"image"`
    DineIn   bool   `json:"dine_in"`
    TakeAway bool   `json:"take_away"`
    Open     bool   `json:"open"`
    AC       bool `json:"ac"`
    NonAC    bool `json:"non_ac"`

	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`

	Token string `json:"token"`
}

var collection = db.Client.Database("eatplek").Collection("restaurant")

func GetRestaurantByUsername(username string) (Restaurant, error) {
	var restaurant Restaurant
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&restaurant)
	if err != nil {
		return restaurant, errors.New("username not found")
	}
	return restaurant, nil
}

func Add(c *gin.Context) (Restaurant, error) {
	var restaurant Restaurant

	c.ShouldBindJSON(&restaurant)

	_, err := GetRestaurantByUsername(restaurant.Username)
	if err == nil {
		return restaurant, errors.New("username already exists")
	}

	restaurant.ID = services.GenerateId()
	restaurant.Open = true
	restaurant.CreatedOn = time.Now()
	restaurant.UpdatedOn = time.Now()
	restaurant.Password, _ = services.HashPassword(restaurant.Password)

	_, err = collection.InsertOne(context.TODO(), restaurant)
	if err != nil {
		return restaurant, err
	}

	return restaurant, nil
}

func Login(c *gin.Context) (LoginResponse, error) {
	var a LoginData
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

	u.Token, err = jwt.GenerateToken(u.ID, "restaurant")
	if err != nil {
		return u, errors.New("error generating token")
	}

	u.Password = "xxxxxxx"
	return u, nil
}

func Get() ([]Restaurant, error) {
	var restaurants []Restaurant

    cur, err := collection.Find(context.TODO(), bson.M{"open":true})
	if err != nil {
		return restaurants, err
	}

	for cur.Next(context.TODO()) {
		var restaurant Restaurant
		err := cur.Decode(&restaurant)
		if err != nil {
			return restaurants, err
		}

		restaurant.Password = "xxxxxxx"
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func GetAll() ([]Restaurant, error) {
	var restaurants []Restaurant

    cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return restaurants, err
	}

	for cur.Next(context.TODO()) {
		var restaurant Restaurant
		err := cur.Decode(&restaurant)
		if err != nil {
			return restaurants, err
		}

		restaurant.Password = "xxxxxxx"
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func GetOne(id string) (Restaurant, error) {
	var restaurant Restaurant

	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&restaurant)
	if err != nil {
		return restaurant, err
	}

	return restaurant, nil
}

func Update(c *gin.Context) error {
	var restaurant Restaurant

	c.ShouldBindJSON(&restaurant)
	restaurant.UpdatedOn = time.Now()

	update := bson.D{{"$set", bson.D{
		{"name", restaurant.Name},
		{"location", restaurant.Location},
		{"phone", restaurant.Phone},
		{"username", restaurant.Username},
		{"email", restaurant.Email},
		{"daysopen", restaurant.DaysOpen},
		{"opentime", restaurant.OpenTime},
		{"closetime", restaurant.CloseTime},
		{"type", restaurant.Type},
		{"image", restaurant.Image},
		{"dinein", restaurant.DineIn},
		{"takeaway", restaurant.TakeAway},
		{"isveg", restaurant.IsVeg},
		{"ac", restaurant.AC},
		{"nonac", restaurant.NonAC},
		{"capacity", restaurant.Capacity},
		{"updatedon", restaurant.UpdatedOn},
	}}}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": c.Param("id")}, update)
	if err != nil {
		return err
	}

	return nil
}

func Delete(id string) error {
	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return err
	}

	_ = food.DeleteMany(id)

	return nil
}

func GetDetails(id string) (RestaurantForOrder, error) {
	var restaurant RestaurantForOrder

	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&restaurant)
	if err != nil {
		return restaurant, err
	}

	return restaurant, nil
}

func UpdateOpen(id string) error {
	type open struct {
		Open bool `bson:"open"`
	}

	var o open
	opts := options.FindOne().SetProjection(bson.D{{"open", 1}})
	err := collection.FindOne(context.TODO(), bson.M{"id": id}, opts).Decode(&o)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.D{{"$set", bson.D{{"open", !o.Open}}}})

	return nil
}

func GetOpenStatus(id string) (bool,error){
	type open struct{
		Open bool `bson:"open"`
	}

	var o open
	opts := options.FindOne().SetProjection(bson.D{{"open", 1}})
	err := collection.FindOne(context.TODO(), bson.M{"id": id}, opts).Decode(&o)
	if err != nil {
		return false,err
	}

	return o.Open,nil
}

func ChangePassword(c *gin.Context) error {
	type password struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var p password
	c.ShouldBindJSON(&p)

	token := c.Request.Header["Token"][0]
	id := jwt.GetUserID(token)

	var restaurant Restaurant
	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&restaurant)
	if err != nil {
		return err
	}

	if !services.CheckPasswordHash(p.OldPassword, restaurant.Password) {
		return errors.New("wrong password")
	}

	p.NewPassword, _ = services.HashPassword(p.NewPassword)

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"id": id},
		bson.D{{"$set", bson.D{{"password", p.NewPassword}}}},
	)
	if err != nil {
		return err
	}

	return nil
}

func ResetPassword(c *gin.Context) error {
	type EmailID struct {
		Email string `json:"email"`
	}

	var email EmailID
	c.ShouldBindJSON(&email)

	var restaurant Restaurant
	err := collection.FindOne(context.TODO(), bson.M{"email": email.Email}).Decode(&restaurant)
	if err != nil {
		return errors.New("email id is not registered")
	}

	jwtToken, err := jwt.GenerateTokenForPasswordReset(restaurant.ID)
	if err != nil {
		return errors.New("error generating token")
	}

	err = services.SendPasswordResetMail(restaurant.Email, restaurant.Name, jwtToken)
	if err != nil {
		return errors.New("error sending mail")
	}

	return nil
}

func ResetAndChangePassword(c *gin.Context) error {
	type resetPassword struct {
		Password string `json:"password"`
	}

	var rp resetPassword
	c.ShouldBindJSON(&rp)

	id := jwt.GetUserID(c.Request.Header["Token"][0])

	rp.Password, _ = services.HashPassword(rp.Password)
	err := collection.FindOneAndUpdate(context.TODO(), bson.M{"id": id}, bson.D{{"$set", bson.D{{"password", rp.Password}}}}).
		Err()
	if err != nil {
		return errors.New("error updating password")
	}

	return nil
}

func Profile(id string) (HotelProfile, error) {
	var p HotelProfile

	findOptions := options.FindOne().SetProjection(bson.M{
		"name":     1,
		"location": 1,
		"phone":    1,
		"email":    1,
	})

	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}, findOptions).Decode(&p)
	if err != nil {
		return HotelProfile{}, errors.New("failed to fetch hotel details")
	}

	return p, nil
}

func UpdateProfile(c *gin.Context) (HotelProfile, error) {
	var p HotelProfile
	c.ShouldBindJSON(&p)

	token := c.Request.Header["Token"][0]
	id := jwt.GetUserID(token)

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"id": id},
		bson.D{{"$set", bson.D{
			{"name", p.Name},
			{"location", p.Location},
			{"phone", p.Phone},
			{"email", p.Email},
		}}},
	)
	if err != nil {
		return HotelProfile{}, err
	}

	return p, nil
}

func GetTimings(id string) (HotelTimings, error) {
	var t HotelTimings

	findOptions := options.FindOne().SetProjection(bson.M{
		"daysopen":  1,
		"opentime":  1,
		"closetime": 1,
        "open": 1,
	})

	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}, findOptions).Decode(&t)
	if err != nil {
		return HotelTimings{}, errors.New("failed to fetch hotel timings")
	}

	return t, nil
}

func UpdateTimings(c *gin.Context) (HotelTimings, error) {
	var t HotelTimings
	c.ShouldBindJSON(&t)

	token := c.Request.Header["Token"][0]
	id := jwt.GetUserID(token)

	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"id": id},
		bson.D{{"$set", bson.D{
			{"daysopen", t.DaysOpen},
			{"opentime", t.OpenTime},
			{"closetime", t.CloseTime},
			{"open",t.Open},
			{"updatedon", time.Now()},
		}}},
	)
	if err != nil {
		return HotelTimings{}, err
	}

	return t, nil
}

func UpdateDeviceToken(c *gin.Context) error {
	type token struct{
		DeviceToken string `json:"device_token"`
		Type string `json:"type"`
	} 

	var t token
	c.ShouldBindJSON(&t)

	auth_token := c.Request.Header["Token"][0]
	id := jwt.GetUserID(auth_token)

	var fieldname string
	if t.Type == "mobile"{
		fieldname = "devicetoken"
	}else{		
		fieldname = "webdevicetoken"
	}

    if t.DeviceToken == ""{
        return errors.New("empty string")
    }

	_,err := collection.UpdateOne(context.TODO(),bson.D{{"id",id}}, bson.D{
        {"$addToSet",bson.D{{fieldname,t.DeviceToken},}},
        {"$set",bson.D{{"updatedon",time.Now()}},},
    })

	if err != nil{
		return err
	}

	return nil
}

func GetDeviceToken(rid string) ([]string,[]string,error){
	type token struct{
        DeviceToken []string `json:"devicetoken" bson:"devicetoken,omitempty"`
		WebDeviceToken []string `json:"webdevicetoken" bson:"webdevicetoken,omitempty"`
	}

	var t token
	opts := options.FindOne().SetProjection(bson.D{{"devicetoken",1},{"webdevicetoken",1}})
	err := collection.FindOne(context.TODO(),bson.D{{"id",rid},},opts).Decode(&t)
	if err != nil {
		return []string{""},[]string{""}, err
	}

    return t.DeviceToken,t.WebDeviceToken, nil
}

func GetOpeningAndClosingTime(rid string)(string,string,error){
    
    type restaurantTiming struct{
        OpenTime string `bson:"opentime,omitempty"`
		CloseTime string `bson:"closetime,omitempty"`
	}

	var t restaurantTiming
	opts := options.FindOne().SetProjection(bson.D{{"opentime",1},{"closetime",1}})
	err := collection.FindOne(context.TODO(),bson.D{{"id",rid},},opts).Decode(&t)
	if err != nil {
		return "","",err
	}

    return t.OpenTime,t.CloseTime, nil

}
