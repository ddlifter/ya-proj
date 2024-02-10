package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	server "1/cmd/app/server"
)

func main() {
	connStr := "user=postgres password=12345 dbname=postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Expressions (id SERIAL PRIMARY KEY, MathExpr TEXT, Status TEXT, Result INTEGER)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/expressions", server.GetExpressions(db)).Methods("GET")
	router.HandleFunc("/api/go/expressions", server.AddExpression(db)).Methods("POST")
	router.HandleFunc("/api/go/expressions/{id}", server.GetExpression(db)).Methods("GET")
	router.HandleFunc("/api/go/expressions/{id}", server.DeleteExpression(db)).Methods("DELETE")

	// wrap the router with CORS and JSON content type middlewares
	enhancedRouter := server.EnableCORS(server.JsonContentTypeMiddleware(router))

	// start server
	fmt.Println("Server is running on :8000")
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))

}
