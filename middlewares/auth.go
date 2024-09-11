package middlewares

import (
	"net/http"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
)

func UsersPreviliges(userType string, expectedUserType string) bool {
	users := map[string]int{
		"user":       1,
		"admin":      2,
		"superAdmin": 3,
	}
	return users[userType] >= users[expectedUserType]
}

func JWTMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// token := c.GetHeader("Bearer-Token")
		aesAccess, _ := c.Cookie("aes-meal-access")
		c.Cookie("aes-meal-refresh")

		userInfo, err := services.VerifyToken(aesAccess, db.TokenTypeAccess)

		if err != nil {
			models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !UsersPreviliges(userInfo.Role, role) {
			models.SendErrorResponse(c, http.StatusUnauthorized, "User is not previlized")
			return
		}
		c.Set("userInfo", userInfo)
		c.Next()
	}
}
