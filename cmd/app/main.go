package main

import (
	agent "1/cmd/app/agent"
	server "1/cmd/app/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/go/", server.Home)
	router.HandleFunc("/api/go/expressions/agent", agent.CalculateHandler).Methods("GET")
	router.HandleFunc("/api/go/expressions", server.GetExpressions).Methods("GET")
	router.HandleFunc("/api/go/expressions", server.AddExpression).Methods("POST")
	router.HandleFunc("/api/go/expression/{id}", server.GetExpression).Methods("GET")
	router.HandleFunc("/api/go/expression/{id}", server.DeleteExpression).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
	// // create router

	// // start server
	// fmt.Println("Server is running on :8000")
	// http.ListenAndServe(":8000", nil)
}
