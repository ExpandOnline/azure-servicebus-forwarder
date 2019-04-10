package main

import (
	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/streadway/amqp"
)

type Publisher struct {
	URL      string
	Exchange string
	output   *Output
	channel  *amqp.Channel
	conn     *amqp.Connection
}

func (pub *Publisher) SetURL(url string) {
	pub.URL = url
}

func (pub *Publisher) SetExchange(exchange string) {
	pub.Exchange = exchange
}

// SetOutput ...
func (pub *Publisher) SetOutput(output *Output) {
	pub.output = output
}

func (pub *Publisher) Close() {
	pub.output.Info.Println("Closing publishing channel")
	pub.channel.Close()
	pub.conn.Close()
}

func (pub *Publisher) Init() {
	conn, err := amqp.Dial(pub.URL)
	if err != nil {
		pub.output.Info.Fatalln(err)
		return
	}
	pub.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		pub.output.Info.Fatalln(err)
		return
	}

	err = channel.ExchangeDeclare(pub.Exchange, "fanout", true, false, false, false, amqp.Table{})
	if err != nil {
		pub.output.Info.Fatalln(err)
		return
	}

	pub.channel = channel
}

func (pub Publisher) Publish(msg *servicebus.Message) error {
	pub.output.Info.Println("Forwarding message with ID " + msg.ID)

	err := pub.channel.Publish(pub.Exchange, "", false, false, pub.transformMessage(msg))
	if err != nil {
		return err
	}

	return nil
}

func (pub Publisher) transformMessage(msg *servicebus.Message) amqp.Publishing {
	return amqp.Publishing{
		ContentType:   msg.ContentType,
		CorrelationId: msg.CorrelationID,
		Body:          msg.Data,
		ReplyTo:       msg.ReplyTo,
		MessageId:     msg.ID,
	}
}
