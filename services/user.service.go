package services

import (
	"errors"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/utils"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func  CreateUser(name string, email string, plainPassword string, employeeId string) (*db.User, error) {
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

// UpdateNote updates a note with id
func UpdateUsersWeeklyMealPlan(userId primitive.ObjectID, request *models.WeeklyMealPlanRequest) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find user")
	}

	if utils.IsTimeIsLessThanGivenTime(9) {
		user.WeeklyMealPlan = request.WeeklyMealPlan
	} else {
		user.PendingWeeklyMealPlan = request.WeeklyMealPlan
	}

	err = mgm.Coll(user).Update(user)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func PendingUsersWeeklyMealPlanService() (*[]db.User, error) {
	users := &[]db.User{}
	userColl := &db.User{}

	filter := bson.M{
		"$where": "this.pendingWeeklyMealPlan.length > 9",
	}

	err := mgm.Coll(userColl).SimpleFind(users, filter)

	if err != nil {
		return nil, errors.New("no pending weekly plan")
	}

	return users, nil
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
