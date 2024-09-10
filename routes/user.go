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
	}
}

func UserRoute(router *gin.RouterGroup) {
	user := router.Group("/user", middlewares.JWTMiddleware("user"))
	{
		user.PUT(
			"/:userId/update-weekly-meal-plan",
			validators.UserWeeklyMealPlanValidator(),
			controllers.UpdateWeeklyMealPlan,
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
	}
}
