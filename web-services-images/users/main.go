package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"web-service-users/src/user/delivery/gin_http"
	"web-service-users/src/user/repository"
	"web-service-users/src/user/usercase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
	log "github.com/sirupsen/logrus"
)

func main() {

	mongodbUser := os.Getenv("MONGODB_USER")
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")
	mongodbHost := os.Getenv("MONGODB_HOST")
	mongodbPort := os.Getenv("MONGODB_PORT")
	webAddr := os.Getenv("WEB_ADDRESS")

	configureLog()

	db := connectDataBase(mongodbUser, mongodbPassword, mongodbHost, mongodbPort)
	userRepo := repository.NewMongodbUserRepository(db)
	userUC := usercase.NewUserUserCase(userRepo)
	engine := createGinHTTPDelivery()
	gin_http.NewUserGinHTTPHandler(engine, userUC)
	engine.Run(webAddr)

}

func connectDataBase(user string, password string, host string, port string) *mongo.Database {

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	client, err := mongo.NewClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("user")
	return db

}

func createGinHTTPDelivery() *gin.Engine {

	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true
	server.Use(cors.New(config))
	return server

}

func configureLog() {
	file, err := os.OpenFile("/var/log/web-service-users/users.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}
