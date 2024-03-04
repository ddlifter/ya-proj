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

// Открытие соединения бд с задачами
func Database() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:12345@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Expressions (id SERIAL PRIMARY KEY, MathExpr TEXT, Status TEXT, Result REAL)")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
