package food

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/services"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Food struct {
	Name        string  `json:"name"         bson:"name"`
	Description string  `json:"description"  bson:"description"`
	NonAcPrice  float64 `json:"non_ac_price,omitempty" bson:"nonacprice,omitempty"`
	AcPrice     float64 `json:"ac_price,omitempty"     bson:"acprice,omitempty"`
	Image       string  `json:"image"        bson:"image"`
	IsVeg       bool    `json:"is_veg"       bson:"isveg"`
	IsAvailable bool    `json:"is_available" bson:"isavailable"`

	RestaurantId   string `json:"restaurant_id"   bson:"restaurantid"`
	RestaurantName string `json:"restaurant_name" bson:"restaurantname"`
	CategoryId     string `json:"category_id"     bson:"categoryid"`
	CategoryName   string `json:"category_name"   bson:"categoryname"`

	ID        string    `json:"id"         bson:"id"`
	CreatedOn time.Time `json:"created_on" bson:"createdon"`
	UpdatedOn time.Time `json:"updated_on" bson:"updatedon"`
}

type FoodsByCategory struct {
	CategoryId   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	Foods        []Food `json:"foods"`
}

type FoodMenu struct {
	ID          string  `json:"id"           bson:"id"`
	Name        string  `json:"name"         bson:"name"`
	NonAcPrice  float64 `json:"non_ac_price" bson:"nonacprice"`
	AcPrice     float64 `json:"ac_price"     bson:"acprice"`
	Description string  `json:"description"  bson:"description"`
}

type FoodForEdit struct{
	ID string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

var collection = db.Client.Database("eatplek").Collection("food")

func Add(c *gin.Context) (Food, error) {
	var food Food

	c.ShouldBindJSON(&food)

	food.IsAvailable = true
	food.ID = services.GenerateId()
	food.CreatedOn = time.Now()
	food.UpdatedOn = time.Now()

	_, err := collection.InsertOne(context.TODO(), food)
	if err != nil {
		return food, err
	}

	return food, nil
}

func Get() ([]Food, error) {
	var foods []Food

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return foods, err
	}

	for cur.Next(context.TODO()) {
		var food Food
		err := cur.Decode(&food)
		if err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func GetOne(id string) (Food, error) {
	var food Food

	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&food)
	if err != nil {
		return food, err
	}

	return food, nil
}

func GetByRestaurant(id,usertype string) ([]FoodsByCategory, error) {
	var foods []FoodsByCategory

	cur, err := collection.Aggregate(context.TODO(), mongo.Pipeline{
		bson.D{{
			"$match", bson.D{
				{"restaurantid", id},
				{"isavailable", true},
			},
		}},
		bson.D{{
			"$group", bson.D{
				{"_id", bson.D{{"category_id", "$categoryid"}, {"category_name", "$categoryname"}}},
				{"items", bson.D{
					{"$push", "$$ROOT"},
				}},
			},
		}},
	})

    if usertype == "admin"{
        cur, err = collection.Aggregate(context.TODO(), mongo.Pipeline{
            bson.D{{
                "$match", bson.D{
                    {"restaurantid", id},
                },
            }},
            bson.D{{
                "$group", bson.D{
                    {"_id", bson.D{{"category_id", "$categoryid"}, {"category_name", "$categoryname"}}},
                    {"items", bson.D{
                        {"$push", "$$ROOT"},
                    }},
                },
            }},
        })
    }

    if err != nil {
        return foods, err
    }

	var response []bson.M
	if err := cur.All(context.TODO(), &response); err != nil {
		return foods, err
	}

	for _, v := range response {
		var food FoodsByCategory

		food.CategoryId = v["_id"].(bson.M)["category_id"].(string)
		food.CategoryName = v["_id"].(bson.M)["category_name"].(string)

		for _, item := range v["items"].(bson.A) {
			var fd Food
			f, _ := bson.Marshal(item)
			_ = bson.Unmarshal(f, &fd)

			food.Foods = append(food.Foods, fd)
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func GetByRestaurantForEdit(id string) ([]FoodForEdit, error){
	var foods []FoodForEdit

	opts := options.Find().SetProjection(bson.D{{"name",1},{"id",1}})
	cur , err := collection.Find(context.TODO(),bson.D{{"restaurantid", id}},opts)
	
	if err != nil {
		return foods, err
	}

	for cur.Next(context.TODO()) {
		var food FoodForEdit
		err := cur.Decode(&food)
		if err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func GetByRestaurantAndCategory(rid, cid string) ([]Food, error) {
	var foods []Food

	cur, err := collection.Find(context.TODO(), bson.M{"restaurantid": rid, "categoryid": cid})
	if err != nil {
		return foods, err
	}

	for cur.Next(context.TODO()) {
		var food Food
		err := cur.Decode(&food)
		if err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func GetByCategory(id string) ([]Food, error) {
	var foods []Food

	cur, err := collection.Find(context.TODO(), bson.M{"categoryid": id})
	if err != nil {
		return foods, err
	}

	for cur.Next(context.TODO()) {
		var food Food
		err := cur.Decode(&food)
		if err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func Update(c *gin.Context) error {
	var food Food

	c.ShouldBindJSON(&food)
	food.UpdatedOn = time.Now()

	update := bson.D{{"$set", bson.D{
		{"name", food.Name},
		{"description", food.Description},
		{"nonacprice", food.NonAcPrice},
		{"acprice", food.AcPrice},
		//{"categoryname", food.CategoryName},
		//{"categoryid", food.CategoryId},
		{"image", food.Image},
		//{"isveg", food.IsVeg},
		{"updatedon", food.UpdatedOn},
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

	return nil
}

func DeleteMany(restaurant_id string) error {
	_, err := collection.DeleteMany(context.TODO(), bson.M{"restaurant_id": restaurant_id})
	if err != nil {
		return err
	}

	return nil
}

func DeleteByCategory(id string) error {
	_, err := collection.DeleteMany(context.TODO(), bson.M{"categoryid": id})
	if err != nil {
		return errors.New("error deleting foods")
	}

	return nil
}

func UpdateCategoryName(id, name string) error {
	update := bson.D{{"$set", bson.D{
		{"categoryname", name},
	}}}

	_, err := collection.UpdateMany(context.TODO(), bson.M{"categoryid": id}, update)
	if err != nil {
		return errors.New("error updating category name")
	}

	return nil
}

func UpdateAvailability(id string) error {
	type available struct {
		Available bool `bson:"isavailable"`
	}

	var a available

	opts := options.FindOne().SetProjection(bson.M{"isavailable": 1})
	err := collection.FindOne(context.TODO(), bson.M{"id": id}, opts).Decode(&a)

    if err != nil{
        return err
    }

    var newval bool
    if a.Available == true{
        newval = false
    }else{
        newval = true
    }
	_, err = collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.D{{"$set", bson.D{
		{"isavailable", newval},
	}}})
	if err != nil {
		return errors.New("error updating availability")
	}

	return nil
}

func MenuChange(id string) ([]FoodMenu, error) {
	var foods []FoodMenu

	cur, err := collection.Find(context.TODO(), bson.M{"restaurantid": id})
	if err != nil {
		return foods, err
	}

	for cur.Next(context.TODO()) {
		var food FoodMenu
		err := cur.Decode(&food)
		if err != nil {
			return foods, err
		}

		foods = append(foods, food)
	}

	return foods, nil
}

func UpdateMenu(c *gin.Context, rid string) (FoodMenu, error) {
	var food FoodMenu
	c.ShouldBindJSON(&food)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": food.ID, "restaurantid": rid}, bson.D{{"$set", bson.D{
		{"name", food.Name},
		{"description", food.Description},
		{"nonacprice", food.NonAcPrice},
		{"acprice", food.AcPrice},
		{"updatedon", time.Now()},
	}}})
	if err != nil {
		return food, errors.New("error updating menu")
	}

	return food, nil
}
