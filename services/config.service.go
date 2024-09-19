package services

import (
	"log"
	"os"

	"github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models"
	"github.com/spf13/viper"
)

var Config *models.EnvConfig

func LoadConfig() {
	v := viper.New()

	// Load the .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		// Load the .env file
		v.SetConfigFile(".env")
		v.SetConfigType("env") // This is important to specify the config type
		if err := v.ReadInConfig(); err != nil {
			log.Printf("Error reading .env file: %s\n", err.Error())
		}
	} else {
		// .env file not found, load environment variables from the system
		log.Println(".env file not found, falling back to environment variables")
		v.AutomaticEnv()

		v.SetDefault("SERVER_ADDR", "0.0.0.0")
		v.SetDefault("SERVER_PORT", "8080")
		v.SetDefault("MONGO_URI", "mongodb://localhost:27017")
		v.SetDefault("MONGO_DATABASE", "exampledb")
		v.SetDefault("USE_REDIS", true)
		v.SetDefault("REDIS_DEFAULT_ADDR", "localhost:6379")
		v.SetDefault("JWT_SECRET", "My.Ultra.Secure.Password")
		v.SetDefault("JWT_ACCESS_EXPIRATION_MINUTES", 1440)
		v.SetDefault("JWT_REFRESH_EXPIRATION_DAYS", 7)
		v.SetDefault("MODE", "debug")
	}

	// Initialize the Config variable
	Config = &models.EnvConfig{}

	// Unmarshal environment variables into the Config struct
	if err := v.Unmarshal(Config); err != nil {
		log.Fatalf("Error unmarshaling config: %s\n", err.Error())
	}

	// Validate the config
	if err := Config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %s\n", err.Error())
	}
}
