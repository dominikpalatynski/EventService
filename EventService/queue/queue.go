package queue

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dominikpalatynski/EventService/storage"
	"github.com/dominikpalatynski/EventService/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	UserId string `json:"user_id"`
	Title string `json:"title"`
}

type QueueHandler struct {
	storage *storage.MongoDbStorage
	context *QueueContext
}

type QueueContext struct {
	connection *amqp.Connection
	channel *amqp.Channel
	exchange string
	routingKey string
}

func NewQueueHandler(s *storage.MongoDbStorage) *QueueHandler{

	ctx, err := newQueueContext()

	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}

	return &QueueHandler{
		storage: s,
		context: ctx,
	}
}

func newQueueContext() (*QueueContext, error) {
	util.LoadEnv()

	exchange := "delayed_exchange"
	routingKey := "delayed_key"
	
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_PORT"))

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	args := amqp.Table{
		"x-delayed-type": "direct",
	}

	err = ch.ExchangeDeclare(
		exchange,     // name
		"x-delayed-message",    // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		args,                   // arguments
	)

	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"hello",
		false,  
		false,   
		false,   
		false,  
		nil,)
	
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &QueueContext{
		connection: conn,
		channel: ch,
		exchange: exchange,
		routingKey: routingKey,
	}, nil
}



func (q *QueueHandler) StartMonitor() {
	ticker := time.NewTicker(1 * time.Minute)

	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	for range ticker.C {
		currentTime := time.Now()
		fmt.Println("fetching events:", currentTime)

		twoMinutesLater := currentTime.Add(2 * time.Minute)

		events, err := q.storage.GetAllEvents(createFilter(currentTime, twoMinutesLater))

		if err != nil {
			fmt.Println("Error fetching events:", err)
            continue
		}

		for _, event := range events {
			delay, err := calculcateDelay(currentTime, event.StartDate)

			if err != nil {
				log.Print("Cannot calculate delay")
				continue
			}
			
			msg, err := createByteMessage(event.UserId, event.Title)

			if err != nil {
				log.Print("Cannot create byte message")
				continue
			}

			if err := q.sendMessage(delay, msg, ctx); err != nil {
				log.Print("Fail during pushing in to queue")
			} else {
				log.Printf("Message sent succesfully at %v with delay %v", currentTime, delay / 1000)
			}
        }
	}
}

func (q *QueueHandler) sendMessage(delay int64, body []byte, ctx context.Context) error {
	err := q.context.channel.PublishWithContext(ctx,
		q.context.exchange, // exchange
		q.context.routingKey,// routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{
				"x-delay": delay,
			},
		})

	return err
}
