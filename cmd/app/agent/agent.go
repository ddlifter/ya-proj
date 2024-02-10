package agent

import (
	"fmt"
	"strconv"
	"strings"
)

func EvaluateExpression(expr string) (float64, error) {
	tokens := tokenize(expr)
	fmt.Println(tokens)
	// Преобразование в обратную польскую нотацию
	rpn := shuntingYard(tokens)
	fmt.Println(rpn)
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
