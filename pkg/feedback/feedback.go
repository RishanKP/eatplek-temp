package feedback

import (
	"context"
	"eatplek/pkg/db"
	"eatplek/pkg/services"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Feedback struct {
	Id        string    `json:"id"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"created_at"`
}

var collection = db.Client.Database("eatplek").Collection("feedback")

func Create(c *gin.Context) error {
	var feedback Feedback

	c.BindJSON(&feedback)
	feedback.Id = services.GenerateId()
	feedback.CreatedAt = time.Now()

	_, err := collection.InsertOne(context.TODO(), feedback)
	if err != nil {
		return errors.New("Error while inserting feedback")
	}

	return nil
}

func Get() ([]Feedback, error) {
	var feedbacks []Feedback

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, errors.New("Error while getting feedbacks")
	}

	for cursor.Next(context.TODO()) {
		var feedback Feedback
		cursor.Decode(&feedback)
		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

func GetOne(id string) (Feedback, error) {
	var feedback Feedback

	err := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&feedback)
	if err != nil {
		return feedback, errors.New("Error while getting feedback")
	}

	return feedback, nil
}

func Delete(id string) error {
	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return errors.New("Error while deleting feedback")
	}

	return nil
}
