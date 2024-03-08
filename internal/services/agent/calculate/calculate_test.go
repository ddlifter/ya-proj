package calculate

import "testing"

func TestEvaluateExpression(t *testing.T) {
	Plus = 0
	Minus = 0
	Mult = 0
	Div = 0
	tests := []struct {
		input    string
		expected float64
	}{
		{"2+2", 4},
		{"(2+2)*3", 12},
		{"(4*2)/2", 4},
		{"(1+2)*(2-1)-(1+1)", 1},
	}

	for _, test := range tests {
		result := EvaluateExpression(test.input)
		if result != test.expected {
			t.Errorf("For expression %s, expected %f but got %f", test.input, test.expected, result)
		}
	}
}
