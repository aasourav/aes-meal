package services

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/utils"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func CreateUser(name string, email string, plainPassword string, employeeId string) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(password), name, db.RoleUser, employeeId)
	user.WeeklyMealPlan = []bool{false, false, false, false, false, false, false}
	err = mgm.Coll(user).Create(user)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

// FindUserById find user by id
func FindUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

func GetUsers() (*[]db.User, error) {
	userDocs := &[]db.User{}
	userCollection := &db.User{}
	err := mgm.Coll(userCollection).SimpleFind(userDocs, bson.M{})
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return userDocs, nil
}

func CreateUpdateUserMeal(user db.User) {
	mealCollection := &db.Meal{}

	dayOfWeek, dayOfMonth, month, year := utils.GetDateDetails()
	err := mgm.Coll(mealCollection).First(bson.M{"consumerId": user.ID, "dayOfWeek": dayOfWeek, "dayOfMonth": dayOfMonth, "year": year, "month": month}, mealCollection)

	if err == nil {
		numberOfMeal := 0
		if user.WeeklyMealPlan[dayOfWeek] {
			numberOfMeal = 1
		}
		mealCollection.NumberOfMeal = numberOfMeal
		err = mgm.Coll(mealCollection).Update(mealCollection)
		if err != nil {
			log.Println("meal update error: ", err.Error())
		}
	} else {
		log.Println("IMPORTANT ERROR: ", err.Error())
		mealCollection := db.NewMeal(user.ID, dayOfWeek, dayOfMonth, month, year)
		err = mgm.Coll(mealCollection).Create(mealCollection)
		if err != nil {
			log.Println("meal create error: ", err.Error())
		}
	}
}

func GetUserMealData(userId primitive.ObjectID, monthNumber string, yearNumber string) (*[]db.Meal, error) {
	mealCollection := &db.Meal{}
	mealDocs := &[]db.Meal{}
	month, _ := strconv.Atoi(monthNumber)
	year, _ := strconv.Atoi(yearNumber)

	err := mgm.Coll(mealCollection).SimpleFind(mealDocs, bson.M{"consumerId": userId, "month": month, "year": year})

	if err != nil {
		return mealDocs, err
	}

	return mealDocs, nil
}

func CleanPendingMeal(userId primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"_id": userId}, user)
	if err != nil {
		return err
	}

	user.PendingWeeklyMealPlan = []bool{}

	err = mgm.Coll(user).Update(user)

	if err != nil {
		return err
	}
	return nil
}

// UpdateNote updates a note with id
func UpdateUsersWeeklyMealPlan(userId primitive.ObjectID, request *models.WeeklyMealPlanRequest) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}
	// fmt.Println("TODAY: ", int(time.Now().Local().Weekday()))
	weekDayIndex := int(time.Now().Local().Weekday())
	// user.WeeklyMealPlan
	if utils.ItTimeIsInRange(12, 22) && user.WeeklyMealPlan[weekDayIndex] == request.WeeklyMealPlan[weekDayIndex] {
		user.WeeklyMealPlan = request.WeeklyMealPlan
	} else {
		user.PendingWeeklyMealPlan = request.WeeklyMealPlan
	}

	err = mgm.Coll(user).Update(user)

	if err != nil {
		return nil, errors.New("cannot update")
	}

	return user, nil
}

func DailyUsersMealData(day int, month int, year int) ([]bson.M, error) {
	// users := &[]db.User{}
	userColl := &db.User{}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "meals"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "consumerId"},
			{Key: "as", Value: "mealsConsumed"},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "mealsConsumed", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$mealsConsumed"},
					{Key: "as", Value: "meal"},
					{Key: "cond", Value: bson.D{
						{Key: "$and", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$$meal.dayOfMonth", day}}},
							bson.D{{Key: "$eq", Value: bson.A{"$$meal.month", month}}},
							bson.D{{Key: "$eq", Value: bson.A{"$$meal.year", year}}},
						}},
					}},
				}},
			}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "mealsConsumed", Value: bson.D{
				{Key: "$ne", Value: bson.A{}},
			}},
		}}},
	}

	cursor, err := mgm.Coll(userColl).Aggregate(context.TODO(), pipeline,
		options.Aggregate().SetMaxTime(time.Minute*1),
		options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M

	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	// if err != nil {
	// 	return nil, errors.New("no pending weekly plan")
	// }

	return results, nil
}

func PendingUsersWeeklyMealPlanService() ([]bson.M, error) {
	// users := &[]db.User{}
	userColl := &db.User{}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "pendingWeeklyMealPlan", Value: bson.D{
				{Key: "$exists", Value: true},
				{Key: "$ne", Value: nil},
				{Key: "$type", Value: "array"},
			}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{
				{Key: "$eq", Value: bson.A{bson.D{{Key: "$size", Value: "$pendingWeeklyMealPlan"}}, 7}},
			}},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "meals"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "consumerId"},
			{Key: "as", Value: "mealCount"},
		}}},
	}

	cursor, err := mgm.Coll(userColl).Aggregate(context.TODO(), pipeline,
		options.Aggregate().SetMaxTime(time.Minute*1),
		options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M

	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	// if err != nil {
	// 	return nil, errors.New("no pending weekly plan")
	// }

	return results, nil
}

func UsersTotalMealByMonthService(month string, year string) ([]bson.M, error) {
	monthInt, _ := strconv.Atoi(month)
	yearInt, _ := strconv.Atoi(year)

	userColl := &db.User{}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "meals"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "consumerId"},
			{Key: "as", Value: "mealConsumption"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$mealConsumption"},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "mealConsumption.month", Value: monthInt},
			{Key: "mealConsumption.year", Value: yearInt},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "totalMeals", Value: bson.D{
				{Key: "$sum", Value: "$mealConsumption.numberOfMeal"},
			}},
			{Key: "name", Value: bson.D{
				{Key: "$first", Value: "$name"},
			}},
			{Key: "employeeId", Value: bson.D{
				{Key: "$first", Value: "$employeeId"},
			}},
		}}},
	}

	cursor, err := mgm.Coll(userColl).Aggregate(context.TODO(), pipeline,
		options.Aggregate().SetMaxTime(time.Minute*1),
		options.Aggregate().SetAllowDiskUse(true))

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M

	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	return results, nil
}

func UpdateUserMealService(mealId string, newMeal string) (*db.Meal, error) {
	mealCollection := &db.Meal{}
	mealObjectId, _ := primitive.ObjectIDFromHex(mealId)
	meal, _ := strconv.Atoi(newMeal)

	err := mgm.Coll(mealCollection).First(bson.M{"_id": mealObjectId}, mealCollection)
	if err != nil {
		return mealCollection, err
	}

	mealCollection.NumberOfMeal = meal

	err = mgm.Coll(mealCollection).Update(mealCollection)
	if err != nil {
		return mealCollection, err
	}

	return mealCollection, nil
}

func ApproveUserWeeklyPlanService(userId string) error {
	user := &db.User{}
	mealCollection := &db.Meal{}
	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	err := mgm.Coll(user).First(bson.M{"_id": userObjectId}, user)
	if err != nil {
		return err
	}

	if len(user.PendingWeeklyMealPlan) > 0 {
		user.WeeklyMealPlan = user.PendingWeeklyMealPlan
		user.PendingWeeklyMealPlan = []bool{}
	}

	err = mgm.Coll(user).Update(user)
	if err != nil {
		return err
	}

	err = mgm.Coll(mealCollection).First(bson.M{"consumerId": userObjectId}, mealCollection)
	if err != nil {
		return err
	}

	dayOfWeek, _, _, _ := utils.GetDateDetails()
	numberOfMeal := mealCollection.NumberOfMeal
	if user.WeeklyMealPlan[dayOfWeek] {
		numberOfMeal = 1
	}
	mealCollection.NumberOfMeal = numberOfMeal

	err = mgm.Coll(mealCollection).Update(mealCollection)
	if err != nil {
		return err
	}
	return nil
}

func RejectUserWeeklyPlanService(userId string) error {
	user := &db.User{}
	userObjectId, _ := primitive.ObjectIDFromHex(userId)
	err := mgm.Coll(user).First(bson.M{"_id": userObjectId}, user)
	if err != nil {
		return err
	}

	user.PendingWeeklyMealPlan = []bool{}

	err = mgm.Coll(user).Update(user)

	if err != nil {
		return err
	}
	return nil
}

// FindUserByEmail find user by email
func FindUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// CheckUserMail search user by email, return error if someone uses
func CheckUserMail(email string) error {
	user := &db.User{}
	userCollection := mgm.Coll(user)
	err := userCollection.First(bson.M{"email": email}, user)
	if err == nil {
		return errors.New("email is already in use")
	}

	return nil
}

// CheckEmployeeId search user by employeeId, return error if someone uses
func CheckEmployeeId(employeeId string) error {
	user := &db.User{}
	userCollection := mgm.Coll(user)
	err := userCollection.First(bson.M{"employeeId": employeeId}, user)
	if err == nil {
		return errors.New("employeeId is already in use")
	}

	return nil
}
