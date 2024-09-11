package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Register
// @Description  registers a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RegisterRequest true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/register [post]

func CronjobAction() {
	users, err := services.GetUsers()
	if err != nil {
		log.Println("cronjob get users error: ", err.Error())
		return
	}

	log.Println("job start")

	for _, user := range *users {
		log.Println("job for: ", user.Name)
		services.CreateUpdateUserMeal(user)
	}
}

func ActionPendingWeeklyPlan(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userId := c.Param("userId")
	actionType := c.Param("actionType")

	if actionType == "approve" {
		err := services.ApproveUserWeeklyPlanService(userId)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}
		response.StatusCode = http.StatusOK
		response.Success = true
	} else if actionType == "reject" {
		err := services.RejectUserWeeklyPlanService(userId)
		if err != nil {
			response.Message = err.Error()
			response.SendResponse(c)
			return
		}
		response.StatusCode = http.StatusOK
		response.Success = true
	}

	response.SendResponse(c)
}

func UpdateUserMeal(c *gin.Context) {
	mealId := c.Param("mealId")
	newMeal := c.Param("newMeal")

	mealData, err := services.UpdateUserMealService(mealId, newMeal)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = map[string]any{
		"mealData": mealData,
	}
	response.SendResponse(c)
}

func PendingWeeklyMealPlans(c *gin.Context) {
	pendingWeeklyPlans, err := services.PendingUsersWeeklyMealPlanService()

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.Data = map[string]any{
		"pendingWeeklyPlans": pendingWeeklyPlans,
	}
	response.SendResponse(c)
}

func CleanPendingMeal(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	userInfo, _ := c.Get("userInfo")
	user, _ := userInfo.(*db.User)

	err := services.CleanPendingMeal(user.ID)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

func UpdateWeeklyMealPlan(c *gin.Context) {
	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// idHex := c.Param("userId")
	// _id, _ := primitive.ObjectIDFromHex(idHex)

	// userId, exists := c.Get("userId")
	// if !exists {
	// 	response.Message = "cannot get user"
	// 	response.SendResponse(c)
	// 	return
	// }
	userInfo, _ := c.Get("userInfo")
	user, _ := userInfo.(*db.User)

	var weeklyMealPlanRequest models.WeeklyMealPlanRequest
	_ = c.ShouldBindBodyWith(&weeklyMealPlanRequest, binding.JSON)

	// userId.(primitive.ObjectID)
	err := services.UpdateUsersWeeklyMealPlan(user.ID, &weeklyMealPlanRequest)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true
	response.SendResponse(c)
}

func Register(c *gin.Context) {
	var requestBody models.RegisterRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// is email in use
	err := services.CheckUserMail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	err = services.CheckEmployeeId(requestBody.EmployeeId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// create user record
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	user, err := services.CreateUser(requestBody.Name, requestBody.Email, requestBody.Password, requestBody.EmployeeId)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// generate access tokens
	// accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	// if err != nil {
	// 	response.Message = err.Error()
	// 	response.SendResponse(c)
	// 	return
	// }

	response.StatusCode = http.StatusCreated
	response.Success = true
	response.Data = gin.H{
		"user": user,
	}
	response.SendResponse(c)

	response.SendResponse(c)
}

// Login godoc
// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.LoginRequest true "Login Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var requestBody models.LoginRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// get user by email
	user, err := services.FindUserByEmail(requestBody.Email)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	// check hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		response.Message = "email and password don't match"
		response.SendResponse(c)
		return
	}

	// generate new access tokens
	accessToken, refreshToken, err := services.GenerateAccessTokens(user)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	response.StatusCode = http.StatusOK
	response.Success = true

	response.Data = gin.H{
		"user": user,
	}
	c.SetCookie("aes-meal-access", *accessToken, 3600, "/", "localhost", true, true)
	c.SetCookie("aes-meal-refresh", *refreshToken, 3600, "/", "localhost", true, true)
	response.SendResponse(c)
}

// Refresh godoc
// @Summary      Refresh
// @Description  refreshes a user token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.RefreshRequest true "Refresh Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Router       /auth/refresh [post]
func Refresh(c *gin.Context) {
	var requestBody models.RefreshRequest
	_ = c.ShouldBindBodyWith(&requestBody, binding.JSON)

	response := &models.Response{
		StatusCode: http.StatusBadRequest,
		Success:    false,
	}

	// check token validity
	token, err := services.VerifyToken(requestBody.Token, db.TokenTypeRefresh)
	if err != nil {
		response.Message = err.Error()
		response.SendResponse(c)
		return
	}

	fmt.Println("refresh :", token)

	// user, err := services.FindUserById("")
	// if err != nil {
	// 	response.Message = err.Error()
	// 	response.SendResponse(c)
	// 	return
	// }

	// // delete old token
	// err = services.DeleteTokenById("token.ID")
	// if err != nil {
	// 	response.Message = err.Error()
	// 	response.SendResponse(c)
	// 	return
	// }

	// accessToken, refreshToken, _ := services.GenerateAccessTokens(user)
	// response.StatusCode = http.StatusOK
	// response.Success = true
	// response.Data = gin.H{
	// 	"user": user,
	// 	// "token": gin.H{
	// 	// 	"access":  accessToken.GetResponseJson(),
	// 	// 	"refresh": refreshToken.GetResponseJson()},
	// }
	// c.SetCookie("aes-meal", *accessToken+"(AES-Meal)"+*refreshToken, 3600, "/", "localhost", true, true)
	response.SendResponse(c)
}
