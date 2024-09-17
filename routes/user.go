package routes

import (
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/controllers"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/middlewares"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/middlewares/validators"
	"github.com/gin-gonic/gin"
)

func UserAuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST(
			"/register",
			validators.RegisterValidator(),
			controllers.Register,
		)

		auth.POST(
			"/login",
			validators.LoginValidator(),
			controllers.Login,
		)

		auth.POST(
			"/refresh",
			validators.RefreshValidator(),
			controllers.Refresh,
		)
		auth.GET(
			"/user",
			controllers.UserAuthorization,
		)
		auth.GET(
			"/logout",
			controllers.UserLogout,
		)
	}
}

func UserRoute(router *gin.RouterGroup) {
	user := router.Group("/user", middlewares.JWTMiddleware("user"))
	{
		user.PUT(
			"/update-weekly-meal-plan",
			validators.UserWeeklyMealPlanValidator(),
			controllers.UpdateWeeklyMealPlan,
		)
		user.DELETE(
			"/clean-pending-meal",
			controllers.CleanPendingMeal,
		)
		user.GET(
			"/user-meal-data/month/:monthNumber/year/:yearNumber",
			controllers.GetMealData,
		)
	}
}

func UserAdminRoute(router *gin.RouterGroup) {
	user := router.Group("/super-user", middlewares.JWTMiddleware("admin"))
	{
		user.GET(
			"/get-pending-weekly-meal-plan",
			controllers.PendingWeeklyMealPlans,
		)
		user.PUT(
			"/action-pending-weekly-meal-plan/action/:actionType/user/:userId",
			controllers.ActionPendingWeeklyPlan,
		)
		user.GET(
			"/meal-data-signeture/day/:day/month/:month/year/:year",
			controllers.UsersDailyMeal,
		)
		user.PUT(
			"/edit-user-meal-plan/meal/:mealId/new-meal/:newMeal",
			controllers.UpdateUserMeal,
		)
		user.GET(
			"/users-total-meal/month/:month/year/:year",
			controllers.UsersTotalMealByMonth,
		)
	}
}
