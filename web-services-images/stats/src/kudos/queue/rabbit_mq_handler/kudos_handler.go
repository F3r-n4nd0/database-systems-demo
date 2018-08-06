package rabbit_mq_handler

import (
	"stats/src/kudos"
	"stats/src/model"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func NewKudosQueueHttpHandler(channel *amqp.Channel, routingKey string, kuc kudos.UseCase) {

	msgs, err := channel.Consume(
		routingKey,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			dataKudos, err := model.UnMarshal(d.Body)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Info("Kudos update sent to rabbit mq")
			kuc.QueueKudos(dataKudos)
		}
	}()

	<-forever
}
