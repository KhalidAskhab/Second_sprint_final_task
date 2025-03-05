package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"SSFT/pkg/calculation"
	"SSFT/pkg/models"
)

var (
	computingPower    int
	internalTaskURL   = "http://localhost:8080/internal/task"   // URL для получения задачи
	internalResultURL = "http://localhost:8080/internal/result" // URL для отправки результата
)

func init() {
	// Чтение переменной среды COMPUTING_POWER
	computingPower = getEnvAsInt("COMPUTING_POWER", 1)
}

func Start() {
	for {
		task, err := getTask()
		if err != nil {
			log.Printf("Ошибка при получении задачи: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		result, err := performCalculation(task)
		if err != nil {
			log.Printf("Ошибка при выполнении вычисления: %v\n", err)
			continue
		}

		if err := sendResult(task.ID, result); err != nil {
			log.Printf("Ошибка при отправке результата: %v\n", err)
		}

		log.Printf("Задача с ID %s успешно обработана, результат: %f\n", task.ID, result)
		time.Sleep(2 * time.Second)
	}
}

func getTask() (models.Task, error) {
	resp, err := http.Get(internalTaskURL) // Используем переменную internalTaskURL
	if err != nil {
		return models.Task{}, fmt.Errorf("ошибка при запросе задачи: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Task{}, fmt.Errorf("оркестратор вернул статус: %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return models.Task{}, fmt.Errorf("ошибка при декодировании задачи: %v", err)
	}

	log.Printf("Получена задача: %+v\n", task)
	return task, nil
}

func performCalculation(task models.Task) (float64, error) {
	// Формируем строку выражения
	var expressionBuilder strings.Builder
	for i, num := range task.Numbers {
		expressionBuilder.WriteString(fmt.Sprintf("%f", num))
		if i < len(task.Operators) {
			expressionBuilder.WriteString(" " + task.Operators[i] + " ")
		}
	}

	expression := expressionBuilder.String()
	result, err := calculation.Calc(expression)
	if err != nil {
		return 0, fmt.Errorf("ошибка при вычислении выражения: %v", err)
	}

	// Корректируем время выполнения в зависимости от COMPUTING_POWER
	time.Sleep(time.Duration(100/computingPower) * time.Millisecond)

	return result, nil
}

func sendResult(taskID string, result float64) error {
	resultData := models.Result{
		ID:     taskID,
		Result: result,
	}

	data, err := json.Marshal(resultData)
	if err != nil {
		return fmt.Errorf("ошибка при кодировании результата: %v", err)
	}

	resp, err := http.Post(internalResultURL, "application/json", bytes.NewBuffer(data)) // Используем переменную internalResultURL
	if err != nil {
		return fmt.Errorf("ошибка при отправке результата: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("оркестратор вернул статус: %d", resp.StatusCode)
	}

	log.Printf("Результат для задачи с ID %s успешно отправлен\n", taskID)
	return nil
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
