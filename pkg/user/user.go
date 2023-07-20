package user

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/jwt"
	"eatplek/pkg/services"
	"errors"
	"fmt"
	"time"
    "math/rand"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
//	Password  string    `json:"password"`
//	Verified  bool      `json:"verified"`
	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserForOrder struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type OTP struct {
	OTP       int       `json:"otp"`
	OTPExpire time.Time `json:"otpexpire"`
	Phone     string    `json:"phone"`
}

type LoginResponse struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	ID        string    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
	Token     string    `json:"token"`
}

var (
	collection    = db.Client.Database("eatplek").Collection("users")
	otpCollection = db.Client.Database("eatplek").Collection("otp")
)

func EmailExists(email string) bool {
	var temp User
	err := collection.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&temp)
	if err != nil {
		return false
	}
	return true
}

func SendOTP(c *gin.Context) error {
	var otp OTP
	c.BindJSON(&otp)

	otp.OTP = rand.Intn(999999-100000)+100000
	otp.OTPExpire = time.Now().Add(time.Minute * 2)

    msg := fmt.Sprintf("Your eatplek verification code is %d.Thank you for choosing eatplek.",otp.OTP) 
    err := services.SendOTP(otp.Phone,msg,"OTP_TID")
	if err != nil {
		return err
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"phone", otp.Phone}}
	update := bson.D{
		{"$set", bson.D{{"otp", otp.OTP}}},
		{"$set", bson.D{{"otpexpire", otp.OTPExpire}}},
	}
	_, err = otpCollection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func Login(c *gin.Context) (LoginResponse, error) {
	var login OTP
	c.BindJSON(&login)

	var temp OTP
	err := otpCollection.FindOne(context.TODO(), bson.D{{"phone", login.Phone}}).Decode(&temp)
	if err != nil {
		return LoginResponse{}, err
	}

	if temp.OTP != login.OTP {
		return LoginResponse{}, errors.New("incorrect otp")
	}
	if temp.OTPExpire.Before(time.Now()) {
		return LoginResponse{}, errors.New("otp expired")
	}

	var user LoginResponse
	err = collection.FindOne(context.TODO(), bson.D{{"phone", login.Phone}}).Decode(&user)
	if err != nil {
		var u User
		u.Phone = login.Phone
		u.ID = services.GenerateId()
		u.CreatedOn = time.Now()

		fmt.Println("creating new user")
		_, err := collection.InsertOne(context.TODO(), u)
		if err != nil {
			return user, errors.New("failed to create new user")
		}

		user.Phone = u.Phone
		user.ID = u.ID
		user.CreatedOn = u.CreatedOn
	}

    go func(){
        _,err = otpCollection.DeleteOne(context.TODO(),bson.D{{"phone",login.Phone}})    
        if err != nil{
            fmt.Println("failed to delete otp record")
        }
    }()

	user.Token, err = jwt.GenerateToken(user.ID, "user")
	if err != nil {
		return user, err
	}

	return user, nil
}

//func Register(c *gin.Context) error {
//	var user User
//	var err error
//	c.BindJSON(&user)
//
//	if EmailExists(user.Email) {
//		return errors.New("email already exists")
//	}
//
//	user.ID = services.GenerateId()
//	user.CreatedOn = time.Now()
//	user.Password, err = services.HashPassword(user.Password)
//	if err != nil {
//		return err
//	}
//
//	user.Verified = false
//
//	_, err = collection.InsertOne(context.Background(), user)
//	if err != nil {
//		return err
//	}
//
//	token, err := jwt.GenerateTokenForAccountVerification(user.ID)
//	if err != nil {
//		return err
//	}
//
//	err = services.SendVerificationEmail(user.Email, user.Name, token)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func Login(c *gin.Context) (LoginResponse, error) {
//	var cr Credentials
//	c.BindJSON(&cr)
//
//	var user User
//	err := collection.FindOne(context.TODO(), bson.D{{"email", cr.Email}}).Decode(&user)
//	if err != nil {
//		return LoginResponse{}, errors.New("invalid email")
//	}
//
//	if !services.CheckPasswordHash(cr.Password, user.Password) {
//		return LoginResponse{}, errors.New("invalid password")
//	}
//
//	if !user.Verified {
//		return LoginResponse{}, errors.New("email not verified")
//	}
//
//	token, err := jwt.GenerateToken(user.ID, "user")
//	if err != nil {
//		return LoginResponse{}, err
//	}
//
//	return LoginResponse{
//		Name:      user.Name,
//		Phone:     user.Phone,
//		Email:     user.Email,
//		ID:        user.ID,
//		CreatedOn: user.CreatedOn,
//		Token:     token,
//	}, nil
//}

func Verify(c *gin.Context) error {
	type token struct {
		Token string `json:"token"`
	}
	var t token

	c.BindJSON(&t)
	if jwt.IsTokenExpired(t.Token) {
		return errors.New("Token Expired")
	}

	id := jwt.GetUserID(t.Token)
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"verified", true},
		}},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func Update(c *gin.Context) error {
	var user User
	c.BindJSON(&user)

	filter := bson.D{{"id", user.ID}}
	update := bson.D{
		{"$set", bson.D{{"name", user.Name}, {"email", user.Email}}},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(id string) (User, error) {
	var temp User
	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}).Decode(&temp)
	if err != nil {
		return temp, err
	}
	return temp, nil
}

func GetDetails(id string) (UserForOrder, error) {
	var temp UserForOrder

	opts := options.FindOne().SetProjection(bson.D{{"name", 1}, {"phone", 1}})

	err := collection.FindOne(context.TODO(), bson.D{{"id", id}}, opts).Decode(&temp)
	if err != nil {
		return temp, errors.New("failed to fetch user details")
	}
	return temp, nil
}

func UpdateUser(c *gin.Context) error {
	type updateUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var user updateUser

	c.BindJSON(&user)

	id := jwt.GetUserID(c.Request.Header["Token"][0])
	filter := bson.D{{"id", id}}
	update := bson.D{
		{"$set", bson.D{{"name", user.Name}, {"email", user.Email}}},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}


