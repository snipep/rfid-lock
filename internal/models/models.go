package models

import (
	"database/sql"
	// "time"

	"github.com/snipep/iot/internal/database"
)

type Log struct {
	ID     string    `json:"id"`
	Status int      `json:"status"`
	Time   string `json:"time"`
}

type Test struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetLogs() (*Log, error) {
	db := database.GetDB() // Use the shared database instance
	query := "SELECT * FROM log WHERE userid = ? LIMIT 1"
	row := db.QueryRow(query, "Vishnu")

	log := &Log{}
	err := row.Scan(&log.ID, &log.Status, &log.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, nil
	}

	return log, nil
}

func GetUser() (*Test, error) {
	db := database.GetDB() // Use the shared database instance
	query := "SELECT id, name FROM test WHERE id = ? LIMIT 1"
	row := db.QueryRow(query, 1)

	user := &Test{}
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
