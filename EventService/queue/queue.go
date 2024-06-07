package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dominikpalatynski/EventService/storage"
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
}

type QueueContext struct {
	connection *amqp.Connection
	channel *amqp.Channel
}

func NewQueueHandler(s *storage.MongoDbStorage) *QueueHandler{
	server := &QueueHandler{
		storage: s,
	}
	return server
}

func (q *QueueHandler) StartMonitor() {
	ticker := time.NewTicker(1 * time.Minute)

	defer ticker.Stop()

	conn, err := amqp.Dial("amqp://admin:password@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	args := amqp.Table{
		"x-delayed-type": "direct",
	}
	err = ch.ExchangeDeclare(
		"delayed_exchange",     // name
		"x-delayed-message",    // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		args,                   // arguments
	)
	failOnError(err, "Failed to declare the delayed exchange")

	queue, err := ch.QueueDeclare(
		"hello",
		false,  
		false,   
		false,   
		false,  
		nil,)
	
	  failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		queue.Name,
		"delayed_key",
		"delayed_exchange",
		false,
		nil,
	)
	failOnError(err, "Failed to bind the queue")

	  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	  defer cancel()
	  
	  body := "Hello World!"


	  failOnError(err, "Failed to publish a message")
	  log.Printf(" [x] Sent %s\n", body)
	for range ticker.C {
		currentTime := time.Now()
		fmt.Println("fetching events:", currentTime)

		delay := int64(15 * 1000) // 2 minuty w milisekundach; możesz zmienić na dynamiczną wartość

		err = ch.PublishWithContext(ctx,
			"delayed_exchange", // exchange
			"delayed_key",      // routing key
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
