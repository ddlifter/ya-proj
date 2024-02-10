package main

import (
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	//agent "1/cmd/app/agent"
	server "1/cmd/app/server"
)

func main() {

	// create router
	http.HandleFunc("/api/go/", server.Home)
	http.HandleFunc("/api/go/getexpressions", server.GetExpressions)
	http.HandleFunc("/api/go/addexpressions", server.AddExpression)
	http.HandleFunc("/api/go/getexpression/{id}", server.GetExpression)
	http.HandleFunc("/api/go/delexpression/{id}", server.DeleteExpression)

	// start server
	fmt.Println("Server is running on :8000")
	http.ListenAndServe(":8000", nil)
}
