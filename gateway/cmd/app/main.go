package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	server "1/internal/handlers"
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

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Адрес вашего фронтенда
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	handler := c.Handler(router)
	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8000", handler))
}
