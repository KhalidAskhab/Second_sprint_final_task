package application

import (
	"SSFT/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	expressionsMutex = &sync.Mutex{}
	expressions      = make(map[string]*models.Expression)
	tasks            = make(chan models.Task, 100)
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := &Config{
		Addr: os.Getenv("PORT"),
	}
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	numbers, operators, err := parseExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expressionID := generateUniqueID()

	expr := &models.Expression{
		ID:         expressionID,
		Expression: req.Expression,
		Status:     "pending",
	}

	expressionsMutex.Lock()
	expressions[expressionID] = expr
	expressionsMutex.Unlock()

	task := models.Task{
		ID:        expressionID,
		Numbers:   numbers,
		Operators: operators,
	}

	select {
	case tasks <- task:
		log.Printf("Задача с ID %s добавлена в очередь", expressionID)
		expressionsMutex.Lock()
		expr.Status = "processing"
		expressionsMutex.Unlock()
	default:
		http.Error(w, "Очередь задач переполнена", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": expressionID})
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()

	expressionList := make([]models.Expression, 0, len(expressions))
	for _, expr := range expressions {
		expressionList = append(expressionList, *expr)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressionList,
	})
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()

	expr, found := expressions[id]
	if !found {
		http.Error(w, "Выражение не найдено", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(expr)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-tasks:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "Нет доступных задач", http.StatusNotFound)
	}
}

func ReceiveResultHandler(w http.ResponseWriter, r *http.Request) {
	var result models.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	expressionsMutex.Lock()
	defer expressionsMutex.Unlock()

	if expr, ok := expressions[result.ID]; ok {
		expr.Status = "completed"
		expr.Result = result.Result
		log.Printf("Updated expression %s: result=%f", result.ID, result.Result)
	}

	w.WriteHeader(http.StatusOK)
}

func parseExpression(expr string) ([]float64, []string, error) {
	// Разделяем выражение на числа и операторы
	parts := strings.Fields(expr)
	if len(parts) < 3 || len(parts)%2 == 0 {
		return nil, nil, fmt.Errorf("неверный формат выражения")
	}

	var numbers []float64
	var operators []string

	for i, part := range parts {
		if i%2 == 0 {
			// Это число
			num, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("неверное число: %s", part)
			}
			numbers = append(numbers, num)
		} else {
			// Это оператор
			if part != "+" && part != "-" && part != "*" && part != "/" {
				return nil, nil, fmt.Errorf("неподдерживаемый оператор: %s", part)
			}
			operators = append(operators, part)
		}
	}

	return numbers, operators, nil
}

func generateUniqueID() string {
	return uuid.New().String()
}

func (a *Application) RunServer() error {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", AddExpressionHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", GetExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", GetExpressionByIDHandler).Methods("GET")
	r.HandleFunc("/internal/task", GetTaskHandler).Methods("GET")
	r.HandleFunc("/internal/result", ReceiveResultHandler).Methods("POST")

	log.Printf("Сервер запущен на порту %s\n", a.config.Addr)
	return http.ListenAndServe(":"+a.config.Addr, r)
}
