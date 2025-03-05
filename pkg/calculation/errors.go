package calculation

import "errors"

var (
	ErrInvalidExpression  = errors.New("некорректное выражение")
	ErrInvalidParentheses = errors.New("несбалансированные скобки")
	ErrInvalidZero        = errors.New("деление на ноль")
	ErrInvalidOperand     = errors.New("неподдерживаемый оператор")
	ErrInvalidValuesCount = errors.New("недостаточно значений для операции")
	ErrInvalidCalculation = errors.New("ошибка вычисления")
)
