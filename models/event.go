package models

import (
	"time"
	"example.com/rest-api/db"
)

type Event struct {
	ID					int64
	Name				string		`binding:"required"`
	Description	string		`binding:"required"`
	Location		string		`binding:"required"`
	DateTime		time.Time	`binding:"required"`
	UserID			int
}

const timeLayout = "2006-01-02 15:04:05-07:00"

var events = []Event{}

func (e Event) Save() error {
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		var dateTime string
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		event.DateTime, err = time.Parse(timeLayout, dateTime)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}