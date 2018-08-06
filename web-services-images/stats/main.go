package main

import (
	"fmt"
	"os"

	"stats/src/kudos/usercase"

	"stats/src/kudos/queue/rabbit_mq_handler"

	"stats/src/kudos/webservice"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {

	rabbitMqHost := os.Getenv("RABBIT_MQ_HOST")
	rabbitMqPort := os.Getenv("RABBIT_MQ_PORT")
	hostWebServiceKudos := os.Getenv("WEB_SERVICE_KUDOS_HOST")
	portWebServiceKudos := os.Getenv("WEB_SERVICE_KUDOS_PORT")
	hostWebServiceUsers := os.Getenv("WEB_SERVICE_USERS_HOST")
	portWebServiceUsers := os.Getenv("WEB_SERVICE_USERS_PORT")

	configureLog()

	rabbitMqChannel, routingKey := configureRabbitMq(rabbitMqHost, rabbitMqPort)
	wsKudos := webservice.NewKudosWebService(hostWebServiceKudos, portWebServiceKudos)
	wsUsers := webservice.NewUsersWebService(hostWebServiceUsers, portWebServiceUsers)
	kudosUC := usercase.NewKudosUserCase(wsUsers, wsKudos)

	forever := make(chan bool)
	rabbit_mq_handler.NewKudosQueueHttpHandler(rabbitMqChannel, routingKey, kudosUC)
	<-forever

}

func configureRabbitMq(rabbitMqHost string, rabbitMqPort string) (*amqp.Channel, string) {

	connectionURL := fmt.Sprintf("amqp://guest:guest@%s:%s/", rabbitMqHost, rabbitMqPort)
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		log.Panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}
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

func configureLog() {
	file, err := os.OpenFile("/var/log/web-service-stats/stats.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}
