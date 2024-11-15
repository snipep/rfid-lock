package handler

import (
	"encoding/json"
	"net/http"
	// "time"

	// "time"

	"github.com/snipep/iot/internal/models"
)

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
