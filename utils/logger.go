package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// EventPayload defines the structure of the event to be sent
type EventPayload struct {
	EventName  string                 `json:"eventName"`
	UserId     string                 `json:"userId"`
	Properties map[string]interface{} `json:"properties"`
}

// LogEventToProducer sends an event to the external producer service
func LogEventToProducer(eventName, userId string, properties map[string]interface{}) {
	payload := EventPayload{
		EventName:  eventName,
		UserId:     userId,
		Properties: properties,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal event payload:", err)
		return
	}

	resp, err := http.Post("http://18.206.244.146:3001/event", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to send event to Producer:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Event sent to Producer, status code:", resp.StatusCode)
}
