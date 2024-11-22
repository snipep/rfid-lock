package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	// "time"

	"github.com/gorilla/mux"
	MQTT "github.com/snipep/iot/internal/MQTT"

	// "time"
	mqttPkg "github.com/eclipse/paho.mqtt.golang"
	"github.com/snipep/iot/internal/models"
)

// ProcessMessage processes incoming MQTT messages
func ProcessMessage(client mqttPkg.Client, msg mqttPkg.Message) {
	go func() { // Start a Goroutine for each message
		fmt.Printf("Message received on topic '%s': %s\n", msg.Topic(), string(msg.Payload()))

		switch msg.Topic() {
		case "rfid/auth":
			User_Authentication(msg)
			client := MQTT.CreateClient("controller", nil)
			MQTT.Publish(client, "rfid/auth/status", 0, false, "msg from backend")
		case "rfid/auth/status":
			ProcessAuthStatus(msg)
		case "register/user/valid":
			Example(msg, true)
		case "register/user/invalid":
			Example(msg, false)
		default:
			fmt.Printf("Unhandled topic: %s\n", msg.Topic())
		}
	}()
}

// InitializeController sets up the MQTT connection and subscription.
func InitializeController() {
	handler := ProcessMessage

	// Create the MQTT client
	client := MQTT.CreateClient("controller", handler)

	// Subscribe to the topic
	MQTT.Subscribe(client, MQTT.Topics, 0)
	// fmt.Printf("Subscribed to topic: %s\n", MQTT.Topic)
}



func GetUserData(w http.ResponseWriter, r *http.Request) {
	// Call the GetUser function to fetch the user
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log, err := models.GetUser(id)
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
		Access int 	`json:"access"`
	}{
		ID:   log.ID,
		Name: log.Name,
		Access: log.Access,
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
	// Extract the user ID from the URL path parameter
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		// Return an error if the ID is invalid
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call GetLogs to fetch all logs for the given user ID
	logs, err := models.GetLogs(id)
	if err != nil {
		// Return an error if there was a problem fetching logs
		http.Error(w, "Error fetching user logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if no logs were found for the user
	if len(logs) == 0 {
		http.Error(w, "No logs found for the user", http.StatusNotFound)
		return
	}

	// Prepare the response structure, which will be a list of logs
	logsResponse := make([]map[string]interface{}, len(logs))
	for i, log := range logs {
		logsResponse[i] = map[string]interface{}{
			"id":     log.ID,
			"time":   log.Time, // Convert time to string
			"status": log.Status,
		}
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal logs data to JSON
	jsonData, err := json.Marshal(logsResponse)
	if err != nil {
		// Return an error if marshalling fails
		http.Error(w, "Error marshaling log data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserInfo

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		fmt.Println("Error reading body:", err)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON to the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Call the InsertUser function for registration
	err = models.InsertUser(user)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		fmt.Println("Error registering user:", err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func User_Authentication(msg mqttPkg.Message) {
	var id struct{ID int}
	err := json.Unmarshal([]byte(msg.Payload()), &id)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
	name, isValid := models.IsValidRFID(id.ID)
	if isValid == 1{
		fmt.Printf("Welcome Name %s, you are authorized\n", name)
		err = models.InsertLog(id.ID, 1)
		if err != nil {
			fmt.Println("Error inserting data:", err)
		}else {
			fmt.Println("Data inserted successfully")
		}
	} else {
		fmt.Printf("%s Access Denied, you aren't authorized\n", name)
		err = models.InsertLog(id.ID, 0)
		if err != nil {
			fmt.Println("Error inserting data:", err)
		}else {
			fmt.Println("Data inserted successfully")
		}
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	var user models.UserInfo

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		fmt.Println("Error reading body:", err)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON to the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Call the InsertUser function for registration
	err = models.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error editing user", http.StatusInternalServerError)
		fmt.Println("Error editing user:", err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Edited successfully"))
}

func ChangeAccess(w http.ResponseWriter, r *http.Request) {
	var user models.UserInfo

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		fmt.Println("Error reading body:", err)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON to the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Call the InsertUser function for registration
	err = models.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error editing authorization", http.StatusInternalServerError)
		fmt.Println("Error editing  authorization:", err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Authorization changed successfully"))	
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = models.Delete(id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		fmt.Println("Error deleting user:", err)
		return
	}

	// Respond to the client
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User deleted successfully"))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		fmt.Println("Error fetching users:", err)
		return
	}


	// Check if no logs were found for the user
	if len(users) == 0 {
		http.Error(w, "No logs found for the user", http.StatusNotFound)
		return
	}

	// Prepare the response structure, which will be a list of logs
	UsersResponse := make([]map[string]interface{}, len(users))
	for i, user := range users {
		UsersResponse[i] = map[string]interface{}{
			"name":     user.Name,
			"status": user.Access,
		}
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal logs data to JSON
	jsonData, err := json.Marshal(UsersResponse)
	if err != nil {
		// Return an error if marshalling fails
		http.Error(w, "Error marshaling log data", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// func UserAuthentication(msg mqttPkg.Message) {
// 	fmt.Println("Authenticating user with payload:", string(msg.Payload()))
// 	// Add your authentication logic here
// }

func ProcessAuthStatus(msg mqttPkg.Message) {
	fmt.Println("Processing auth status with payload:", string(msg.Payload()))
	// Add your status processing logic here
}

func Example(msg mqttPkg.Message, isValid bool) {
	if isValid {
		fmt.Println("Registering valid user:", string(msg.Payload()))
	} else {
		fmt.Println("Invalid user registration attempt:", string(msg.Payload()))
	}
	// Add user registration logic here
}
