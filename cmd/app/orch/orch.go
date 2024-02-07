package orch

import (
	"encoding/json"
	"net/http"
	"time"
)

// Expression структура для представления арифметического выражения
type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result,omitempty"`
}

// Operation структура для представления доступной операции
type Operation struct {
	Name      string        `json:"name"`
	Execution time.Duration `json:"executionTime"`
}

var Expressions []Expression
var operations []Operation

// AddExpressionHandler обработчик для добавления вычисления арифметического выражения
func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var newExpr Expression
	err := json.NewDecoder(r.Body).Decode(&newExpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление выражения в список
	newExpr.Status = "pending"
	Expressions = append(Expressions, newExpr)

	// Отправка ответа
	w.WriteHeader(http.StatusCreated)
}

// GetExpressionsHandler обработчик для получения списка выражений со статусами
func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	// Отправка списка выражений в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Expressions)
}

// GetExpressionByIDHandler обработчик для получения значения выражения по его идентификатору
func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра из URL
	id := r.URL.Query().Get("id")

	// Поиск выражения по ID
	var foundExpr Expression
	for _, expr := range Expressions {
		if expr.ID == id {
			foundExpr = expr
			break
		}
	}

	// Отправка выражения в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundExpr)
}

// GetAvailableOperationsHandler обработчик для получения списка доступных операций со временем их выполения
func GetAvailableOperationsHandler(w http.ResponseWriter, r *http.Request) {
	// Отправка списка операций в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operations)
}

// GetTaskHandler обработчик для получения задачи для выполения
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получение первой задачи из списка
	if len(Expressions) == 0 {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	task := Expressions[0]
	Expressions = Expressions[1:]

	// Отправка задачи в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// SubmitResultHandler обработчик для приёма результата обработки данных
func SubmitResultHandler(w http.ResponseWriter, r *http.Request) {
	var result Expression
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление статуса и результата в соответствующей задаче
	for i, task := range Expressions {
		if task.ID == result.ID {
			Expressions[i].Status = result.Status
			Expressions[i].Result = result.Result
			break
		}
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}
