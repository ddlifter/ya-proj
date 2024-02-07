package main

import (
	"fmt"
)

func main() {
	expression := "(2+2) * (2+2)"
	result, err := EvaluateExpression(expression)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
