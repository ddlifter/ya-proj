package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type expression struct {
	ID         int
	Expression string
	Status     string
	Result     float64
}

type operation struct {
	Name string
	Time time.Duration
}

type orchestrator struct {
	expressions map[int]*expression
	operations  []operation
	taskQueue   chan int
	taskResults chan int
	mutex       sync.RWMutex
}

func newOrchestrator() *orchestrator {
	return &orchestrator{
		expressions: make(map[int]*expression),
		operations: []operation{
			{Name: "addition", Time: time.Second},
			{Name: "subtraction", Time: 2 * time.Second},
			{Name: "multiplication", Time: 3 * time.Second},
			{Name: "division", Time: 4 * time.Second},
		},
		taskQueue:   make(chan int),
		taskResults: make(chan int),
	}
}

func (o *orchestrator) addExpression(expr string) int {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	id := len(o.expressions) + 1
	newExpr := &expression{
		ID:         id,
		Expression: expr,
		Status:     "pending",
		Result:     0,
	}
	o.expressions[id] = newExpr

	o.taskQueue <- id

	return id
}

func (o *orchestrator) getExpressions() []*expression {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	exprs := make([]*expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		exprs = append(exprs, expr)
	}
	return exprs
}

func (o *orchestrator) getExpressionByID(id int) (*expression, bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	expr, ok := o.expressions[id]
	return expr, ok
}

func (o *orchestrator) getOperations() []operation {
	return o.operations
}

func (o *orchestrator) getTask() int {
	return <-o.taskResults
}

func (o *orchestrator) receiveResult(id int) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	expr, ok := o.expressions[id]
	if !ok {
		return
	}
	expr.Status = "completed"
	delete(o.expressions, id)
}

func processExpression(orch *orchestrator, exprID int, expr string) {
	time.Sleep(5 * time.Second) // Simulating actual expression evaluation process
	result := 0.0               // Replace with actual evaluation logic

	orch.mutex.Lock()
	defer orch.mutex.Unlock()

	exprObj, ok := orch.expressions[exprID]
	if !ok {
		return
	}

	exprObj.Status = "in progress"
	exprObj.Result = result
}

func handleExpressions(w http.ResponseWriter, r *http.Request, orch *orchestrator) {
	switch r.Method {
	case http.MethodGet:
		expressions := orch.getExpressions()
		for _, expr := range expressions {
			fmt.Fprintf(w, "Expression ID: %d, Expression: %s, Status: %s, Result: %f\n", expr.ID, expr.Expression, expr.Status, expr.Result)
		}
	case http.MethodPost:
		expression := r.FormValue("expression")
		id := orch.addExpression(expression)
		fmt.Fprintf(w, "Expression ID: %d\n", id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleExpressionByID(w http.ResponseWriter, r *http.Request, orch *orchestrator) {
	switch r.Method {
	case http.MethodGet:
		id := r.FormValue("id")
		exprID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid expression ID", http.StatusBadRequest)
			return
		}
		expr, ok := orch.getExpressionByID(exprID)
		if !ok {
			http.Error(w, "Expression not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Expression ID: %d, Expression: %s, Status: %s, Result: %f\n", expr.ID, expr.Expression, expr.Status, expr.Result)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleOperations(w http.ResponseWriter, r *http.Request, orch *orchestrator) {
	switch r.Method {
	case http.MethodGet:
		operations := orch.getOperations()
		for _, op := range operations {
			fmt.Fprintf(w, "Operation: %s, Execution Time: %s\n", op.Name, op.Time.String())
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
