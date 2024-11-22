package models

import (
	"database/sql"
	"log"
	"time"

	// "time"

	"github.com/snipep/iot/internal/database"
)

type UserInfo struct {
	ID   	int    	`json:"id"`
	Name 	string 	`json:"name"`
	Access 	int 	`json:"access"`
}

type Log struct {
	ID     string    `json:"id"`
	Status int      `json:"status"`
	Time   string `json:"time"`
}

type Test struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetLogs(id int) ([]Log, error) {
	db := database.GetDB() // Use the shared database instance
	query := "SELECT user_id, log, status FROM rfid WHERE user_id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log
	for rows.Next() {		
		var log Log
		err := rows.Scan(&log.ID, &log.Time, &log.Status)
		if err != nil {
			return nil, err
		}
		// log.Time = log.Time.Format("2006-01-02 15:04:05")
		logs = append(logs, log)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func GetUser(id int) (*UserInfo, error) {
	db := database.GetDB() // Use the shared database instance
	query := "SELECT user_id, name, access FROM user WHERE user_id = ? LIMIT 1"
	row := db.QueryRow(query, id)

	user := &UserInfo{}
	err := row.Scan(&user.ID, &user.Name, &user.Access)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}


func InsertLog(ID int, status int) error {
	db := database.GetDB()
	query := "INSERT INTO rfid (user_id, log, status) VALUES (?, ?, ?)"
	now := time.Now()
	dateTime := now.Format("2006-01-02 15:04:05")
	_, err := db.Exec(query, ID, dateTime, status)
	if err != nil {
		return err
	}
	return nil
}

func IsValidRFID(rfid int) (string, int) {
    db := database.GetDB()
    query := "SELECT name, access FROM user WHERE user_id = ? LIMIT 1"
    row := db.QueryRow(query, rfid)

    var Auth struct {
        Name   string
        Access int
    }

    err := row.Scan(&Auth.Name, &Auth.Access)
    if err != nil {
        if err == sql.ErrNoRows {
            return "Empty", 0 
		}
        return "User not Found", 0
    }

    return Auth.Name, Auth.Access
}

func InsertUser(user UserInfo) error {
	db := database.GetDB()
	query := "INSERT INTO user (user_id, name, access) VALUES (?, ?, ?)"

	_, err := db.Exec(query, user.ID, user.Name, user.Access)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(user UserInfo) error {
	db := database.GetDB()
	query := "UPDATE user SET name = ?, access = ? WHERE user_id = ?"

	_, err := db.Exec(query, user.Name, user.Access, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func EditAccess(id int, access int) error {
	db := database.GetDB()
	query := "UPDATE user SET access = ? WHERE user_id = ?"

	_, err := db.Exec(query, access, id)
	if err != nil {
		return err
	}
	return nil
}

func Delete(id int) error {
	db := database.GetDB()
	query := "DELETE FROM user WHERE user_id = ?"

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]UserInfo, error) {
	db := database.GetDB() // Use your shared database instance

	// SQL query to get all users' names and access levels
	query := "SELECT name, access FROM user"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error querying the database:", err)
		return nil, err
	}
	defer rows.Close()

	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		// Scan each row into the User struct
		err := rows.Scan(&user.Name, &user.Access)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Println("Error after iterating over rows:", err)
		return nil, err
	}

	return users, nil
}