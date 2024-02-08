package main

import (
	a "1/cmd/app/orch"
	"fmt"
	"net/http"
)

func main() {
	// // for i := 0; i < b.AgentCount; i++ {
	// // 	b.Agents[i] = b.Agent{ID: i, State: "Waiting"}
	// // }
	http.HandleFunc("/add-expression", a.AddExpressionHandler)
	http.HandleFunc("/get-expressions", a.GetExpressionsHandler)
	http.HandleFunc("/get-expression-by-id", a.GetExpressionByIDHandler)
	// http.HandleFunc("/get-available-operations", a.GetAvailableOperationsHandler)
	// http.HandleFunc("/get-task", a.GetTaskHandler)
	// http.HandleFunc("/submit-result", a.SubmitResultHandler)
	http.HandleFunc("/start", a.StartHandler)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)

}

// {
// 	"id": "2",
// 	"mathExpr": "2+5",
// 	"status": "pending"
// }
