package queue

import (
	"fmt"
	"time"

	"github.com/dominikpalatynski/EventService/storage"
)

const timeFormat string = "2006-01-02T15:04:05Z"

type QueueHandler struct {
	storage *storage.MongoDbStorage
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

	for range ticker.C {
		currentTime := time.Now()
		fmt.Println("fetching events:", currentTime)
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
