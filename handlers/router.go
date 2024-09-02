package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func CreateRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware configuration
	config := cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	router.Use(cors.New(config))

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/", healthCheckV1)
			// v1.GET("/todos", getTodos)
			// v1.GET("/todos/:id", getTodoById)
			// v1.POST("/todos/create", createTodo)
			// v1.PUT("/todos/update/:id", updateTodo)
			// v1.DELETE("/todos/delete/:id", deleteTodo)
		}

		// version 2 - add it if you want
		// v2 := api.Group("/v2")
		// {
		// }
	}

	return router
}

func healthCheckV1(c *gin.Context) {
	res := Response{
		Msg:  "api v1 is working",
		Code: 200,
	}
	c.JSON(http.StatusOK, res)
}
