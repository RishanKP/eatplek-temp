package category

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/food"
	"eatplek/pkg/services"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Category struct {
	Name string `json:"name"`
	ID   string `json:"id"`

	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
}

var collection = db.Client.Database("eatplek").Collection("category")

func Add(c *gin.Context) (Category, error) {
	var category Category

	c.ShouldBindJSON(&category)

	category.ID = services.GenerateId()
	category.CreatedOn = time.Now()
	category.UpdatedOn = time.Now()

	_, err := collection.InsertOne(context.TODO(), category)
	if err != nil {
		return category, err
	}

	return category, nil
}

func Get() ([]Category, error) {
	var categories []Category

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return categories, err
	}

	for cur.Next(context.TODO()) {
		var category Category
		err := cur.Decode(&category)
		if err != nil {
			return categories, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func Update(c *gin.Context) error {
	var category Category

	c.ShouldBindJSON(&category)
	category.UpdatedOn = time.Now()

	update := bson.D{{"$set", bson.D{
		{"name", category.Name},
		{"updated_on", category.UpdatedOn},
	}}}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": c.Param("id")}, update)
	if err != nil {
		return err
	}

	err = food.UpdateCategoryName(c.Param("id"), category.Name)
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

	err = food.DeleteByCategory(id)
	if err != nil {
		return err
	}

	return nil
}
