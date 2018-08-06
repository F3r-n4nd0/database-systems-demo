package main

import (
	"fmt"
	"os"
	"time"
	"web-service-kudos/src/kudos/delivery/gin_http"
	"web-service-kudos/src/kudos/repository"
	"web-service-kudos/src/kudos/usercase"

	"web-service-kudos/src/kudos/queue"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {

	webAddr := os.Getenv("WEB_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rabbitMqHost := os.Getenv("RABBIT_MQ_HOST")
	rabbitMqPort := os.Getenv("RABBIT_MQ_PORT")

	configureLog()

	clientDb := connectDataBase(redisHost, redisPort, redisPassword)
	rabbitMqChannel, routingKey := configureRabbitMq(rabbitMqHost, rabbitMqPort)

	kudosRepo := repository.NewRedisKudosRepository(clientDb)
	kudosQueue := queue.NewRabbitMqKudosQueue(rabbitMqChannel, routingKey)
	kudosUC := usercase.NewKudosUsecase(kudosRepo, kudosQueue)
	engine := createGinHTTPDelivery()
	gin_http.NewKudosGinHTTPHandler(engine, kudosUC)
	engine.Run(webAddr)

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
	file, err := os.OpenFile("/var/log/web-service-kudos/kudos.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

func connectDataBase(redisHost string, redisPort string, redisPassword string) *redis.Client {

	address := fmt.Sprintf("%s:%s", redisHost, redisPort)
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

}

func configureRabbitMq(rabbitMqHost string, rabbitMqPort string) (*amqp.Channel, string) {

	connectionURL := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitMqHost, rabbitMqPort)
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		log.Panic(err)
	}
	//defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}
	//defer ch.Close()

	routingKey := "kudos"
	q, err := ch.QueueDeclare(
		routingKey,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}
	return ch, q.Name
}
