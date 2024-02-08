package orch

import (
	a "1/cmd/app/agent"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Expression структура для представления арифметического выражения

// Operation структура для представления доступной операции
type Operation struct {
	Name      string        `json:"name"`
	Execution time.Duration `json:"executionTime"`
}

var Expressions = make(map[string]a.Expression)
var operations []Operation

func StartHandler(w http.ResponseWriter, r *http.Request) {
	tasks := make(chan a.Expression, 100) // Буферизированный канал для задач
	ctx, Cancel := context.WithCancel(context.Background())
	for i := 0; i < a.AgentCount; i++ {
		a.Wg.Add(1)
		go a.Worker(ctx, &a.Wg, tasks) // Запускаем горутины
	}

	for _, expr := range Expressions {
		tasks <- expr // Отправляем задачи в канал
	}
	close(tasks) // Закрываем канал после отправки всех задач

	a.Wg.Wait() // Ждем завершения всех горутин
	Cancel()    // Отменяем контекст
}

// AddExpressionHandler обработчик для добавления вычисления арифметического выражения
func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var newExpr a.Expression
	err := json.NewDecoder(r.Body).Decode(&newExpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавление выражения в список
	Expressions[newExpr.ID] = newExpr

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
	var foundExpr a.Expression
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

	//task := Expressions[0]
	//Expressions = Expressions[1:]

	// Отправка задачи в формате JSON
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(task)
}

// SubmitResultHandler обработчик для приёма результата обработки данных
func SubmitResultHandler(w http.ResponseWriter, r *http.Request) {
	var result a.Expression
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление статуса и результата в соответствующей задаче
	for _, task := range Expressions {
		if task.ID == result.ID {
			//Expressions[i].Status = result.Status
			//Expressions[i].Result = result.Result
			break
		}
	}

	// Отправка ответа
	w.WriteHeader(http.StatusOK)
}
