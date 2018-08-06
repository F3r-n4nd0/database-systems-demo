package queue

import (
	"web-service-kudos/src/kudos"

	"web-service-kudos/src/model"

	"encoding/json"

	"github.com/streadway/amqp"
)

type rabbitMqQueue struct {
	Channel    *amqp.Channel
	RoutingKey string
}

func NewRabbitMqKudosQueue(channel *amqp.Channel, routingKey string) kudos.Queue {
	return &rabbitMqQueue{channel, routingKey}
}

type kudosQueue struct {
	UserName string `json:"user_name"`
	Type     string `json:"type"`
}

func (kq *kudosQueue) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(kq)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (r *rabbitMqQueue) IncreasesKudos(kudos *model.Kudos) error {

	data := kudosQueue{
		UserName: kudos.ToUserName,
		Type:     "ADD",
	}
	body, err := data.Marshal()
	if err != nil {
		return err
	}
	return r.send(body)

}

func (r *rabbitMqQueue) DecreasesKudos(kudos *model.Kudos) error {

	data := kudosQueue{
		UserName: kudos.ToUserName,
		Type:     "REMOVE",
	}
	body, err := data.Marshal()
	if err != nil {
		return err
	}
	return r.send(body)

}

func (r *rabbitMqQueue) send(body []byte) error {

	err := r.Channel.Publish(
		"",
		r.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil

}
