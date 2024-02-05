package main

import (
	"1/calculating"
	"1/handling"
	"fmt"
	"net/http"
)

func main() {
	expression := "(2+2) * (2+2)"
	result, err := calculating.EvaluateExpression(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// Регистрация обработчика HTTP-запросов
	http.HandleFunc("/calculate", handling.ExpressionHandler)
	// Запуск HTTP-сервера на порту 8080
	http.ListenAndServe(":8080", nil)
}
