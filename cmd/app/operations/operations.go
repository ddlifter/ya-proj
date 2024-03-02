package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	calculate "1/cmd/app/calculate"
	database "1/cmd/app/database"
)

// Функция для обновления времени выполнения операций
func UpdateOperations(w http.ResponseWriter, r *http.Request) {
	db := database.DbOperation()
	defer db.Close()
	var ops database.Operations
	err := json.NewDecoder(r.Body).Decode(&ops)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление времени операций в бд
	_, err = db.Exec("UPDATE operations SET time = $1 WHERE operation = $2", ops.Time, ops.Operation)
	if err != nil {
		log.Fatal(err)
	}

	// Получение операций для обновления в приложении
	rows, err := db.Query("SELECT * FROM operations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Обновление измененных операций
	for rows.Next() {
		var operation, time string
		if err := rows.Scan(&operation, &time); err != nil {
			log.Fatal(err)
		}
		switch operation {
		case "plus":
			calculate.Plus, _ = strconv.Atoi(time)
		case "minus":
			calculate.Minus, _ = strconv.Atoi(time)
		case "mult":
			calculate.Mult, _ = strconv.Atoi(time)
		case "div":
			calculate.Div, _ = strconv.Atoi(time)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
