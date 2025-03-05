package calculation

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
	}{
		{"Simple addition", "1 + 2", 3, false},
		{"Simple subtraction", "5 - 3", 2, false},
		{"Simple multiplication", "2 * 3", 6, false},
		{"Simple division", "6 / 2", 3, false},
		{"Complex expression", "2 + 3 * 4", 14, false},
		{"Division by zero", "1 / 0", 0, true},
		{"Invalid expression", "1 +", 0, true},
		{"Parentheses", "(1 + 2) * 3", 9, false},
		{"Invalid parentheses", "(1 + 2", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calc(tt.expression)
			if tt.expectError {
				if err == nil {
					t.Errorf("Ожидалась ошибка для выражения: %s", tt.expression)
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка для выражения: %s: %v", tt.expression, err)
				}
				if result != tt.expected {
					t.Errorf("Ожидаемый результат: %f, получено: %f для выражения: %s", tt.expected, result, tt.expression)
				}
			}
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
	}{
		{"Simple addition", "1+2", 3, false},
		{"Simple subtraction", "5-3", 2, false},
		{"Simple multiplication", "2*3", 6, false},
		{"Simple division", "6/2", 3, false},
		{"Complex expression", "2+3*4", 14, false},
		{"Division by zero", "1/0", 0, true},
		{"Invalid expression", "1+", 0, true},
		{"Parentheses", "(1+2)*3", 9, false},
		{"Invalid parentheses", "(1+2", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluateexpression(tt.expression)
			if tt.expectError {
				if err == nil {
					t.Errorf("Ожидалась ошибка для выражения: %s", tt.expression)
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка для выражения: %s: %v", tt.expression, err)
				}
				if result != tt.expected {
					t.Errorf("Ожидаемый результат: %f, получено: %f для выражения: %s", tt.expected, result, tt.expression)
				}
			}
		})
	}
}

func TestAttachOperator(t *testing.T) {
	tests := []struct {
		name        string
		op          rune
		values      []float64
		expected    []float64
		expectError bool
	}{
		{"Addition", '+', []float64{2, 3}, []float64{5}, false},
		{"Subtraction", '-', []float64{5, 3}, []float64{2}, false},
		{"Multiplication", '*', []float64{2, 3}, []float64{6}, false},
		{"Division", '/', []float64{6, 2}, []float64{3}, false},
		{"Division by zero", '/', []float64{1, 0}, []float64{}, true},
		{"Invalid operand", 'x', []float64{1, 2}, []float64{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := attachOperator(tt.op, tt.values)
			if tt.expectError {
				if err == nil {
					t.Errorf("Ожидалась ошибка для операции: %c", tt.op)
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка для операции: %c: %v", tt.op, err)
				}
				if len(result) != len(tt.expected) || (len(result) > 0 && result[0] != tt.expected[0]) {
					t.Errorf("Ожидаемый результат: %v, получено: %v для операции: %c", tt.expected, result, tt.op)
				}
			}
		})
	}
}

func TestSearchNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		start    int
		expected float64
		next     int
	}{
		{"Single digit", "1+2", 0, 1, 1},
		{"Multiple digits", "123+456", 0, 123, 3},
		{"Decimal number", "1.23+4.56", 0, 1.23, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, next := searchnumbers(tt.input, tt.start)
			if result != tt.expected || next != tt.next {
				t.Errorf("Ожидаемый результат: %f, %d; получено: %f, %d для ввода: %s", tt.expected, tt.next, result, next, tt.input)
			}
		})
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		op       rune
		expected int
	}{
		{"Addition", '+', 1},
		{"Subtraction", '-', 1},
		{"Multiplication", '*', 2},
		{"Division", '/', 2},
		{"Unknown operator", 'x', 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := precedence(tt.op)
			if result != tt.expected {
				t.Errorf("Ожидаемый приоритет: %d, получено: %d для оператора: %c", tt.expected, result, tt.op)
			}
		})
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"No empty strings", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"With empty strings", []string{"a", "", "b", "", "c"}, []string{"a", "b", "c"}},
		{"All empty strings", []string{"", "", ""}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeEmptyStrings(tt.input)
			if !equalSlices(result, tt.expected) {
				t.Errorf("Ожидаемый результат: %v, получено: %v для ввода: %v", tt.expected, result, tt.input)
			}
		})
	}
}

// Вспомогательная функция для сравнения слайсов
func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
