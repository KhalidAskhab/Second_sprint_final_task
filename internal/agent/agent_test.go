package agent

import (
	"SSFT/pkg/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetTask(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := models.Task{
			ID:        "123",
			Numbers:   []float64{1, 2, 3},
			Operators: []string{"+", "-"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}))
	defer server.Close()

	// Подменяем URL на тестовый сервер
	oldURL := internalTaskURL
	internalTaskURL = server.URL
	defer func() { internalTaskURL = oldURL }()

	task, err := getTask()
	if err != nil {
		t.Fatalf("Ошибка при получении задачи: %v", err)
	}

	if task.ID != "123" {
		t.Errorf("Ожидаемый ID задачи: 123, получено: %s", task.ID)
	}

	if len(task.Numbers) != 3 || task.Numbers[0] != 1 || task.Numbers[1] != 2 || task.Numbers[2] != 3 {
		t.Errorf("Ожидаемые числа: [1, 2, 3], получено: %v", task.Numbers)
	}

	if len(task.Operators) != 2 || task.Operators[0] != "+" || task.Operators[1] != "-" {
		t.Errorf("Ожидаемые операторы: [+, -], получено: %v", task.Operators)
	}
}

func TestPerformCalculation(t *testing.T) {
	task := models.Task{
		Numbers:   []float64{1, 2, 3},
		Operators: []string{"+", "-"},
	}

	result, err := performCalculation(task)
	if err != nil {
		t.Fatalf("Ошибка при выполнении вычисления: %v", err)
	}

	// Ожидаемый результат зависит от логики функции calculation.Calc
	expectedResult := 0.0 // Замените на ожидаемый результат
	if result != expectedResult {
		t.Errorf("Ожидаемый результат: %f, получено: %f", expectedResult, result)
	}
}

func TestSendResult(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resultData models.Result
		err := json.NewDecoder(r.Body).Decode(&resultData)
		if err != nil {
			t.Fatalf("Ошибка при декодировании результата: %v", err)
		}

		if resultData.ID != "123" {
			t.Errorf("Ожидаемый ID результата: 123, получено: %s", resultData.ID)
		}

		if resultData.Result != 42.0 {
			t.Errorf("Ожидаемый результат: 42.0, получено: %f", resultData.Result)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Подменяем URL на тестовый сервер
	oldURL := internalResultURL
	internalResultURL = server.URL
	defer func() { internalResultURL = oldURL }()

	err := sendResult("123", 42.0)
	if err != nil {
		t.Fatalf("Ошибка при отправке результата: %v", err)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Устанавливаем переменную окружения
	os.Setenv("TEST_ENV", "42")
	defer os.Unsetenv("TEST_ENV")

	// Проверяем, что функция возвращает правильное значение
	value := getEnvAsInt("TEST_ENV", 1)
	if value != 42 {
		t.Errorf("Ожидаемое значение: 42, получено: %d", value)
	}

	// Проверяем значение по умолчанию
	value = getEnvAsInt("NON_EXISTENT_ENV", 1)
	if value != 1 {
		t.Errorf("Ожидаемое значение по умолчанию: 1, получено: %d", value)
	}
}
