package transport

import (
	"fin_quotes/internal/config"
	"github.com/streadway/amqp"
	"log"
)

type Rabbitmq struct {
	Chan  *amqp.Channel
	Queue amqp.Queue
}

func New() *Rabbitmq {
	return &Rabbitmq{}
}

func (rabbit *Rabbitmq) InitConn(cfg *config.Config) {
	// Установите соединение с RabbitMQ
	conn, err := amqp.Dial(cfg.GetRabbitDSN())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// Создайте канал
	ch, err := conn.Channel()
	rabbit.Chan = ch

	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
}

func (rabbit *Rabbitmq) ConnClose() {
	rabbit.Chan.Close()
}

func (rabbit *Rabbitmq) DeclareQueue(name string) {
	args := amqp.Table{
		"x-message-ttl": int32(60000), // TTL 60 секунд
	}

	queue, err := rabbit.Chan.QueueDeclare(
		name,  // имя очереди
		false, // durable постоянная очередь
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // аргументы
	)
	rabbit.Queue = queue
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
}

func (rabbit *Rabbitmq) SendMsg(data []byte) {
	err := rabbit.Chan.Publish(
		"",                // обменник
		rabbit.Queue.Name, // ключ маршрутизации (имя очереди)
		false,             // обязательное
		false,             // немедленное
		amqp.Publishing{
			DeliveryMode: amqp.Transient, // сохранять сообщение
			ContentType:  "text/plain",
			Body:         data,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}
