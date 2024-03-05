package main

import (
	server "1/internal/services/gateway"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	log.Print("server has started")

	// Создание маршрутизатора
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/api/go/home", server.Home)
	router.HandleFunc("/api/go/expressions", server.GetExpressions).Methods("GET")
	router.HandleFunc("/api/go/expressions", server.AddExpression).Methods("POST")
	router.HandleFunc("/api/go/expression/{id}", server.GetExpression).Methods("GET")
	router.HandleFunc("/api/go/expression/{id}", server.DeleteExpression).Methods("DELETE")
	router.HandleFunc("/api/go/expressions", server.DeleteExpressions).Methods("DELETE")
	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8000", router))
}
