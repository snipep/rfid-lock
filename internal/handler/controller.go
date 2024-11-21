package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	MQTT "github.com/snipep/iot/internal/MQTT"

	// "time"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/snipep/iot/internal/models"
)

// ProcessMessage handles incoming MQTT messages.
func ProcessMessage(client mqtt.Client, msg mqtt.Message) {

	// fmt.Printf("Received message on topic: %s\n", msg.Topic())
	fmt.Printf("Message payload: %s\n", string(msg.Payload()))

	// Parse JSON payload if needed (temperature and humidity)
	// Example JSON: {"Temp":25.4,"Hum":60.3}
	var message models.Log
	err := json.Unmarshal([]byte(msg.Payload()), &message)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	models.InsertLog(message.ID, message.Status)
	fmt.Println("Data inserted successfully")

}

// InitializeController sets up the MQTT connection and subscription.
func InitializeController() {
	handler := ProcessMessage

	// Create the MQTT client
	client := MQTT.CreateClient("controller", handler)

	// Subscribe to the topic
	MQTT.Subscribe(client, MQTT.Topic, 0)
	// fmt.Printf("Subscribed to topic: %s\n", MQTT.Topic)
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}


func GetData(w http.ResponseWriter, r *http.Request) {
	// Call the GetUser function to fetch the user
	log, err := models.GetUser()
	if err != nil {
		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if no user was found
	if log == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare user data for response
	user := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		ID:   log.ID,
		Name: log.Name,
	}

	w.Header().Set("Content-Type", "application/json")

	// Marshal user data to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshaling user data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonData)
}

func GetUserLog(w http.ResponseWriter, r *http.Request) {
	// Call the GetUser function to fetch the user
	log, err := models.GetLogs()
	if err != nil {
		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if no user was found
	if log == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prepare user data for response
	user := struct {
		ID   	string	`json:"id"`
		Status 	int 	`json:"status"`
		Time 	string	`json:"time"`
	}{
		ID:   log.ID,
		Status: log.Status,
		Time: log.Time,
	}

	w.Header().Set("Content-Type", "application/json")

	// Marshal user data to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshaling user data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonData)
}

// func InsertData(w http.ResponseWriter, r *http.Request)  {
	
// }