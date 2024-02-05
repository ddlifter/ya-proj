package handle

import (
	"fmt"
	"net/http"
	"time"
)

//Реализация сервера

// Функция для обработки арифметического выражения
func processExpression(expression string, duration time.Duration) {
	// Здесь будет логика обработки выражения
	time.Sleep(duration) // Моделируем "долгие" вычисления
	fmt.Println("Expression", expression, "has been processed")
}

// Обработчик HTTP-запросов
func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	expression := r.URL.Query().Get("expression")
	// Запустить обработку выражения в отдельной горутине
	go processExpression(expression, 5*time.Second)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Expression has been accepted for processing"))
}
