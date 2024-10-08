package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/controllers"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/routes"
	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/services"
	cronJob "github.com/robfig/cron/v3"
)

// @title GoLang Rest API Starter Doc
// @version 1.0
// @description GoLang - Gin - RESTful - MongoDB - Redis
// @termsOfService https://swagger.io/terms/

// @contact.name Ebubekir Yiğit
// @contact.url https://github.com/ebubekiryigit
// @contact.email ebubekiryigit6@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer-Token
func main() {
	services.LoadConfig()
	services.InitMongoDB()

	// if services.Config.UseRedis {
	// 	services.CheckRedisConnection()
	// }

	routes.InitGin()
	router := routes.New()

	server := &http.Server{
		Addr:         services.Config.ServerAddr + ":" + services.Config.ServerPort,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	c := cronJob.New()
	_, err := c.AddFunc("*/1 * * * *", controllers.CronjobAction)
	if err != nil {
		log.Panicln("Cronjob failed")
	}
	c.Start()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 15 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	c.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
