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


func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }

const timeFormat string = "2006-01-02T15:04:05Z"

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
	  
	body := "Hello World!"

	for range ticker.C {
		currentTime := time.Now()
		fmt.Println("fetching events:", currentTime)

		delay := int64(15 * 1000) // 2 minuty w milisekundach; możesz zmienić na dynamiczną wartość

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
		failOnError(err, "Failed to publish a delayed message")
		twoMinutesLater := currentTime.Add(2 * time.Minute)


		filterData := map[string]interface{}{
			"start_date": currentTime.Format(timeFormat),
			"end_date":   twoMinutesLater.Format(timeFormat),
		}

		events, err := q.storage.GetAllEvents(filterData)

		if err != nil {
			fmt.Println("Error fetching events:", err)
            continue
		}

		for _, event := range events {
            fmt.Printf("Event Title Title: %s, currentTime: %v", event.Title, filterData["start_date"])

			//ToDo
			//Add rabbitMq sending message to notifservice
        }
	}
}
