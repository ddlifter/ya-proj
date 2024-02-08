package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Expression struct {
	ID       string  `json:"id"`
	MathExpr string  `json:"mathExpr"`
	Status   string  `json:"status"`
	Result   float64 `json:"result,omitempty"`
}

func EvaluateExpression(expr string) (float64, error) {
	tokens := tokenize(expr)
	// Преобразование в обратную польскую нотацию
	rpn := shuntingYard(tokens)
	// Вычисление результата
	result, err := evaluateRPN(rpn)
	return result, err
}

func tokenize(expr string) []string {
	// Разбиваем выражение на токены
	tokens := []string{}
	buffer := ""
	for _, char := range expr {
		if char == ' ' {
			continue
		} else if strings.Contains("+-*/()", string(char)) {
			if len(buffer) > 0 {
				tokens = append(tokens, buffer)
				buffer = ""
			}
			tokens = append(tokens, string(char))
		} else {
			buffer += string(char)
		}
	}
	if len(buffer) > 0 {
		tokens = append(tokens, buffer)
	}
	return tokens
}

func shuntingYard(tokens []string) []string {
	// Алгоритм преобразования в обратную польскую нотацию
	var rpn []string
	var stack []string
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}
	for _, token := range tokens {
		switch {
		case token == "(":
			stack = append(stack, token)
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Удаление "("
		case precedence[token] > 0:
			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
				rpn = append(rpn, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		default: // Операнд
			rpn = append(rpn, token)
		}
	}
	for len(stack) > 0 {
		rpn = append(rpn, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return rpn
}

func evaluateRPN(tokens []string) (float64, error) {
	var stack []float64
	for _, token := range tokens {
		if !strings.Contains("+-*/", token) {
			value, err := strconv.Atoi(token)
			if err != nil {
				return 0, err
			}
			stack = append(stack, float64(value))
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("Invalid expression")
			}
			operand2, operand1 := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				stack = append(stack, operand1+operand2)
			case "-":
				stack = append(stack, operand1-operand2)
			case "*":
				stack = append(stack, operand1*operand2)
			case "/":
				stack = append(stack, operand1/operand2)
			}
		}
	}
	if len(stack) == 1 {
		return stack[0], nil
	}
	return 0, fmt.Errorf("Invalid expression")
}

var Wg sync.WaitGroup

func Worker(id string) {
	defer Wg.Done()
	expr := getExpressionFromServer(id)

	result, _ := EvaluateExpression(expr)
	sendResultToServer(result, id, expr)
}

func getExpressionFromServer(id string) string {
	resp, err := http.Get("http://localhost:8080/get-expression-by-id?id=" + id)
	if err != nil {
		log.Println("Error getting expression from server:", err)
		return ""
	}
	defer resp.Body.Close()

	var calc Expression
	if err := json.NewDecoder(resp.Body).Decode(&calc); err != nil {
		log.Println("Error decoding response:", err)
		return ""
	}

	return calc.MathExpr
}

func sendResultToServer(result float64, workerID string, expr string) {
	calc := Expression{
		ID:       workerID,
		Result:   result,
		Status:   "completed",
		MathExpr: expr,
	}
	jsonData, err := json.Marshal(calc)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/add-expression", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending result to server:", err)
		return
	}
	defer resp.Body.Close()
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Calculating..."))
}
