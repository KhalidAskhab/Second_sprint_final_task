package calculation

import (
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	additionTime       time.Duration
	subtractionTime    time.Duration
	multiplicationTime time.Duration
	divisionTime       time.Duration
)

func init() {
	additionTime = getEnvAsDuration("TIME_ADDITION_MS", 100)
	subtractionTime = getEnvAsDuration("TIME_SUBTRACTION_MS", 200)
	multiplicationTime = getEnvAsDuration("TIME_MULTIPLICATIONS_MS", 300)
	divisionTime = getEnvAsDuration("TIME_DIVISIONS_MS", 400)
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	// Проверка на корректность выражения
	numbers := strings.FieldsFunc(expression, func(r rune) bool {
		return r == '+' || r == '-' || r == '*' || r == '/' || r == '(' || r == ')'
	})
	operators := strings.FieldsFunc(expression, func(r rune) bool {
		return (r >= '0' && r <= '9') || r == '.' || r == '(' || r == ')'
	})

	operators = removeEmptyStrings(operators)

	if len(numbers) != len(operators)+1 {
		return 0, ErrInvalidExpression
	}

	result, err := evaluateexpression(expression)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Вспомогательная функция для удаления пустых строк из слайса
func removeEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func searchnumbers(expression string, index int) (float64, int) {
	start := index
	for index < len(expression) && (isDigit(expression[index]) || expression[index] == '.') {
		index++
	}
	val, _ := strconv.ParseFloat(expression[start:index], 64)
	return val, index
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isOperator(char byte) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func evaluateexpression(expression string) (float64, error) {
	var ops []rune
	var values []float64

	for i := 0; i < len(expression); {
		char := expression[i]
		if isDigit(char) || char == '.' {
			val, nextIndex := searchnumbers(expression, i)
			values = append(values, val)
			i = nextIndex
		} else if char == '(' {
			ops = append(ops, '(')
			i++
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				var err error
				values, err = attachOperator(ops[len(ops)-1], values)
				if err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
			}
			if len(ops) == 0 {
				return 0, ErrInvalidParentheses
			}
			ops = ops[:len(ops)-1]
			i++
		} else if isOperator(char) {
			currentOp := rune(char)
			for len(ops) > 0 && precedence(currentOp) <= precedence(ops[len(ops)-1]) {
				var err error
				values, err = attachOperator(ops[len(ops)-1], values)
				if err != nil {
					return 0, err
				}
				ops = ops[:len(ops)-1]
			}
			ops = append(ops, currentOp)
			i++
		} else {
			return 0, ErrInvalidCalculation
		}
	}

	for len(ops) > 0 {
		op := ops[len(ops)-1]
		if op == '(' {
			return 0, ErrInvalidParentheses
		}
		var err error
		values, err = attachOperator(op, values)
		if err != nil {
			return 0, err
		}
		ops = ops[:len(ops)-1]
	}

	if len(values) != 1 {
		return 0, ErrInvalidValuesCount
	}
	return values[0], nil
}
func getEnvAsDuration(key string, defaultValue int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(defaultValue) * time.Millisecond
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(defaultValue) * time.Millisecond
	}
	return time.Duration(intValue) * time.Millisecond
}

func attachOperator(op rune, values []float64) ([]float64, error) {
	if len(values) < 2 {
		return values, ErrInvalidValuesCount
	}
	a := values[len(values)-1]
	b := values[len(values)-2]
	values = values[:len(values)-2]

	// Добавляем задержку в зависимости от операции
	switch op {
	case '+':
		time.Sleep(additionTime)
		result := b + a
		return append(values, result), nil
	case '-':
		time.Sleep(subtractionTime)
		result := b - a
		return append(values, result), nil
	case '*':
		time.Sleep(multiplicationTime)
		result := b * a
		return append(values, result), nil
	case '/':
		time.Sleep(divisionTime)
		if a == 0 {
			return values, ErrInvalidZero
		}
		result := b / a
		return append(values, result), nil
	default:
		return values, ErrInvalidOperand
	}
}
