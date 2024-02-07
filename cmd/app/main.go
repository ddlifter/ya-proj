package main

import (
	a "1/cmd/app/orch"
	"fmt"
	"net/http"
)

func main() {
	a.Expressions = append(a.Expressions, a.Expression{ID: "1", MathExpr: "2+2", Status: "pending", Result: 0})
	http.HandleFunc("/add-expression", a.AddExpressionHandler)
	http.HandleFunc("/get-expressions", a.GetExpressionsHandler)
	http.HandleFunc("/get-expression-by-id", a.GetExpressionByIDHandler)
	http.HandleFunc("/get-available-operations", a.GetAvailableOperationsHandler)
	http.HandleFunc("/get-task", a.GetTaskHandler)
	http.HandleFunc("/submit-result", a.SubmitResultHandler)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
