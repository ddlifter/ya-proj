package agent

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	calculate "1/cmd/app/calculate"
	database "1/cmd/app/database"
	workers "1/cmd/app/workers"

	_ "github.com/mattn/go-sqlite3"
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
	_, err = db.Exec("UPDATE operations SET time = ? WHERE operation = ?", ops.Time, ops.Operation)
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

// Функция запуска вычислений
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	db := database.Database()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM Expressions WHERE Status='waiting'") //Берем только те задачи, которые еще не были запущены
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Устанавливаем статус выполнения
	var res []database.Expression
	for rows.Next() {
		var u database.Expression
		if err := rows.Scan(&u.ID, &u.MathExpr, &u.Status, &u.Result); err != nil {
			log.Fatal(err)
		}
		res = append(res, u)
		_, err = db.Exec("UPDATE Expressions SET Status='pending' WHERE ID=$1", u.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	workers.Process(res) // Запускаем вычисления

	for _, expr := range res {
		tokens := calculate.Tokenize(expr.MathExpr)
		rpnTokens := calculate.ShuntingYard(tokens)
		result, err := calculate.EvaluateRPN(rpnTokens)
		if err != nil {
			log.Printf("Error evaluating expression: %v", err)
			continue
		}

		// Обновляем бд выполненными задачами
		_, err = db.Exec("UPDATE Expressions SET Status='complete', Result=$1 WHERE MathExpr=$2", result, expr.MathExpr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Функция для вывода агентов
func AgentHandler(w http.ResponseWriter, r *http.Request) {
	db := database.DbAgent()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM agents")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	agents := map[int]workers.Worker{}
	for rows.Next() {
		var u workers.Worker
		if err := rows.Scan(&u.Id, &u.Status, &u.LastPing); err != nil {
			log.Fatal(err)
		}
		agents[u.Id] = u
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(agents)

}
