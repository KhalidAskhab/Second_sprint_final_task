package application

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Second_sprint_final_task/pkg/models"
)

func TestAddExpressionHandler(t *testing.T) {
	// Создаем тестовый запрос с выражением
	reqBody := `{"expression": "1 + 2 - 3"}`
	req := httptest.NewRequest("POST", "/api/v1/calculate", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	AddExpressionHandler(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидаемый статус код: %d, получено: %d", http.StatusCreated, rr.Code)
	}

	// Декодируем ответ
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверяем, что ID выражения возвращен
	if _, ok := response["id"]; !ok {
		t.Error("Ожидаемый ключ 'id' отсутствует в ответе")
	}
}

func TestGetExpressionsHandler(t *testing.T) {
	// Добавляем тестовое выражение
	expressionsMutex.Lock()
	expressions["test-id"] = &models.Expression{
		ID:         "test-id",
		Expression: "1 + 2",
		Status:     "pending",
	}
	expressionsMutex.Unlock()

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	GetExpressionsHandler(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, получено: %d", http.StatusOK, rr.Code)
	}

	// Декодируем ответ
	var response map[string][]models.Expression
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверяем, что тестовое выражение присутствует в ответе
	found := false
	for _, expr := range response["expressions"] {
		if expr.ID == "test-id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Тестовое выражение не найдено в ответе")
	}
}

func TestGetExpressionByIDHandler(t *testing.T) {
	// Добавляем тестовое выражение
	expressionsMutex.Lock()
	expressions["test-id"] = &models.Expression{
		ID:         "test-id",
		Expression: "1 + 2",
		Status:     "pending",
	}
	expressionsMutex.Unlock()

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/api/v1/expressions/test-id", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "test-id"})
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	GetExpressionByIDHandler(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, получено: %d", http.StatusOK, rr.Code)
	}

	// Декодируем ответ
	var expr models.Expression
	if err := json.NewDecoder(rr.Body).Decode(&expr); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверяем, что возвращено правильное выражение
	if expr.ID != "test-id" {
		t.Errorf("Ожидаемый ID выражения: test-id, получено: %s", expr.ID)
	}
}

func TestGetTaskHandler(t *testing.T) {
	// Добавляем тестовую задачу в канал
	task := models.Task{
		ID:        generateUniqueID(), // Используем динамический ID
		Numbers:   []float64{1, 2},
		Operators: []string{"+"},
	}
	tasks <- task

	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/internal/task", nil)
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	GetTaskHandler(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, получено: %d", http.StatusOK, rr.Code)
	}

	// Декодируем ответ
	var returnedTask models.Task
	if err := json.NewDecoder(rr.Body).Decode(&returnedTask); err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверяем, что возвращена правильная задача (без учета ID)
	if !equalSlices(returnedTask.Numbers, task.Numbers) {
		t.Errorf("Ожидаемые числа: %v, получено: %v", task.Numbers, returnedTask.Numbers)
	}
	if !equalSlices(returnedTask.Operators, task.Operators) {
		t.Errorf("Ожидаемые операторы: %v, получено: %v", task.Operators, returnedTask.Operators)
	}
}

func TestReceiveResultHandler(t *testing.T) {
	// Добавляем тестовое выражение
	expressionsMutex.Lock()
	expressions["test-id"] = &models.Expression{
		ID:         "test-id",
		Expression: "1 + 2",
		Status:     "processing",
	}
	expressionsMutex.Unlock()

	// Создаем тестовый запрос с результатом
	result := models.Result{
		ID:     "test-id",
		Result: 3.0,
	}
	reqBody, _ := json.Marshal(result)
	req := httptest.NewRequest("POST", "/internal/result", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	ReceiveResultHandler(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидаемый статус код: %d, получено: %d", http.StatusOK, rr.Code)
	}

	// Проверяем, что статус выражения обновлен
	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()
	if expr, ok := expressions["test-id"]; ok {
		if expr.Status != "completed" {
			t.Errorf("Ожидаемый статус выражения: completed, получено: %s", expr.Status)
		}
		if expr.Result != 3.0 {
			t.Errorf("Ожидаемый результат: 3.0, получено: %f", expr.Result)
		}
	} else {
		t.Error("Выражение не найдено")
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input     string
		numbers   []float64
		operators []string
		err       bool
	}{
		{"1 + 2 - 3", []float64{1, 2, 3}, []string{"+", "-"}, false},
		{"1 * 2 / 3", []float64{1, 2, 3}, []string{"*", "/"}, false},
		{"1 +", nil, nil, true},     // Неверный формат
		{"1 + 2 *", nil, nil, true}, // Неверный формат
	}

	for _, tt := range tests {
		numbers, operators, err := parseExpression(tt.input)
		if tt.err {
			if err == nil {
				t.Errorf("Ожидалась ошибка для ввода: %s", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("Неожиданная ошибка для ввода: %s: %v", tt.input, err)
			}
			if !equalSlices(numbers, tt.numbers) {
				t.Errorf("Ожидаемые числа: %v, получено: %v для ввода: %s", tt.numbers, numbers, tt.input)
			}
			if !equalSlices(operators, tt.operators) {
				t.Errorf("Ожидаемые операторы: %v, получено: %v для ввода: %s", tt.operators, operators, tt.input)
			}
		}
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
