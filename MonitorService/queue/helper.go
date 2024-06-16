package queue

import (
	"encoding/json"
	"log"
	"time"
)

const timeFormat string = "2006-01-02T15:04:05Z"

func createFilter(current time.Time, target time.Time) map[string]interface{}{
	return map[string]interface{}{
		"start_date": current.Format(timeFormat),
		"end_date":   target.Format(timeFormat),
	}
}

func createByteMessage(userId string, title string) ([]byte, error) {
	msg := Message{
		UserId: userId,
		Title: title,
	}

	body, err := json.Marshal(msg)

	if err != nil {
		return nil, err
	}
	
	return body, nil
}

func calculcateDelay(currentTime time.Time, startTimeStr string) (int64, error) {

	startTime, err := time.Parse(timeFormat, startTimeStr)
	
	log.Printf("startValue: %v", startTime)

	if err != nil {
		return 0, err
	}

    currentTimeUTC := currentTime.Add(1 * time.Minute)
    startTimeUTC := startTime

	delayedTime := startTimeUTC.Sub(currentTimeUTC)

	log.Printf("currentTime.Add(1*time.Minute): %v", currentTimeUTC)
	log.Printf("delayedTime %v", delayedTime)

	return delayedTime.Milliseconds() * 1000, nil
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }