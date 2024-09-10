package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
)

const (
	RoleUser = "user"
)

type User struct {
	mgm.DefaultModel      `bson:",inline"`
	Email                 string `json:"email" bson:"email"`
	Password              string `json:"-" bson:"password"`
	EmployeeId            string `json:"employeeId" bson:"employeeId"`
	WeeklyMealPlan        []bool `json:"weeklyMealPlan" bson:"weeklyMealPlan"`
	PendingWeeklyMealPlan []bool `json:"pendingWeeklyMealPlan" bson:"pendingWeeklyMealPlan"`
	Name                  string `json:"name" bson:"name"`
	Role                  string `json:"role" bson:"role"`
	MailVerified          bool   `json:"mail_verified" bson:"mail_verified"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserInfo User `json:"userInfo"`
}

func NewUser(email string, password string, name string, role string, employeeId string) *User {
	return &User{
		Email:        email,
		Password:     password,
		Name:         name,
		EmployeeId:   employeeId,
		Role:         role,
		MailVerified: false,
	}
}

func (model *User) CollectionName() string {
	return "users"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
