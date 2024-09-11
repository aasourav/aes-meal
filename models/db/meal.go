package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meal struct {
	// dayOfWeek, dayOfMonth
	mgm.DefaultModel `bson:",inline"`
	ConsumerId       primitive.ObjectID `json:"consumerId" bson:"consumerId"`
	DayOfMonth       int                `json:"dayOfMonth" bson:"dayOfMonth"`
	DayOfWeek        int                `json:"dayOfWeek" bson:"dayOfWeek"`
	Month            int                `json:"month" bson:"month"`
	Year             int                `json:"year" bson:"year"`
	NumberOfMeal     int                `json:"numberOfMeal" bson:"numberOfMeal"`
}

func NewMeal(consumerId primitive.ObjectID, dayOfWeek int, dayOfMonth int, month int, year int) *Meal {
	return &Meal{
		ConsumerId: consumerId,
		DayOfWeek:  dayOfWeek,
		DayOfMonth: dayOfMonth,
		Month:      month,
		Year:       year,
	}
}

func (model *Meal) CollectionName() string {
	return "meals"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
