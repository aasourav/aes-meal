package middlewares

import (
	"fmt"
	"net/http"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token := c.GetHeader("Bearer-Token")
		aesAccess, _ := c.Cookie("aes-meal-access")
		aesRefresh, _ := c.Cookie("aes-meal-refresh")

		fmt.Println("TOKEN MODEL 1:", aesAccess)
		fmt.Println("TOKEN MODEL 2:", aesRefresh)

		tokenModel, err := services.VerifyToken(aesAccess, db.TokenTypeAccess)
		if err != nil {
			models.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("userIdHex", tokenModel.User.Hex())
		c.Set("userId", tokenModel.ID)

		c.Next()
	}
}
