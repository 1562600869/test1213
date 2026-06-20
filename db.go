package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Hall struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Theme        string `json:"theme"`
	MaxCapacity  int    `json:"max_capacity"`
	Status       string `json:"status"`
}

type Guide struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Language string `json:"language"`
	Status   string `json:"status"`
}

type Reservation struct {
	ID           int64  `json:"id"`
	GuestName    string `json:"guest_name"`
	GuestPhone   string `json:"guest_phone"`
	HallID       int64  `json:"hall_id"`
	HallName     string `json:"hall_name,omitempty"`
	TimeSlot     string `json:"time_slot"`
	PeopleCount  int    `json:"people_count"`
	GuideID      *int64 `json:"guide_id,omitempty"`
	GuideName    string `json:"guide_name,omitempty"`
	CreatedAt    string `json:"created_at"`
}

type ThemeStat struct {
	Theme string `json:"theme"`
	Count int    `json:"count"`
}

func initDB(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS halls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		theme TEXT NOT NULL,
		max_capacity INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT '开放'
	);

	CREATE TABLE IF NOT EXISTS guides (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nickname TEXT NOT NULL,
		phone TEXT NOT NULL,
		language TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT '在职'
	);

	CREATE TABLE IF NOT EXISTS reservations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guest_name TEXT NOT NULL,
		guest_phone TEXT NOT NULL,
		hall_id INTEGER NOT NULL,
		time_slot TEXT NOT NULL,
		people_count INTEGER NOT NULL,
		guide_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (hall_id) REFERENCES halls(id),
		FOREIGN KEY (guide_id) REFERENCES guides(id)
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}
