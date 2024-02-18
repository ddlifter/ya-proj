package database

import (
	"database/sql"
	"log"
)

// Описание сущности задачи
type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result"`
}

// Описание сущности операции
type Operations struct {
	Operation string `json:"operation"`
	Time      string `json:"time"`
}

// Открытие соединения бд с задачами
func Database() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Expressions (id INTEGER PRIMARY KEY AUTOINCREMENT, MathExpr TEXT, Status TEXT, Result REAL)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Открытие соединения бд с операциями
func DbOperation() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS operations (Operation TEXT, Time TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")

	return db
}

// Открытие соединения бд с агентами
func DbAgent() *sql.DB {
	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS agents (id INTEGER PRIMARY KEY, Status TEXT, LastPing INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("PRAGMA journal_mode=WAL")

	return db
}
