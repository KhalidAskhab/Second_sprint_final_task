package orchestrator

import (
	"SSFT/pkg/models"
	"testing"
	"time"
)

func TestDistributeTasks(t *testing.T) {
	// Создаем Orchestrator
	orch := &Orchestrator{
		tasks:   make(chan models.Task, 10),
		results: make(chan models.Result, 10),
		agents:  make(map[string]AgentInfo),
	}

	// Добавляем агентов
	orch.mu.Lock()
	orch.agents["agent1"] = AgentInfo{ID: "agent1", ComputingPower: 10, LastSeen: time.Now()}
	orch.agents["agent2"] = AgentInfo{ID: "agent2", ComputingPower: 20, LastSeen: time.Now()}
	orch.mu.Unlock()

	// Запускаем распределение задач в отдельной горутине
	go orch.distributeTasks()

	// Добавляем задачи в канал
	tasks := []models.Task{
		{ID: "task1", Numbers: []float64{1, 2}, Operators: []string{"+"}},
		{ID: "task2", Numbers: []float64{3, 4}, Operators: []string{"-"}},
	}

	for _, task := range tasks {
		orch.tasks <- task
	}

	// Даем время для обработки задач
	time.Sleep(1 * time.Second)

	// Проверяем, что задачи были распределены
	orch.mu.Lock()
	defer orch.mu.Unlock()

	// Проверяем, что агент с наибольшей вычислительной мощностью получил задачи
	if len(orch.agents) != 2 {
		t.Errorf("Ожидалось 2 агента, получено: %d", len(orch.agents))
	}

	// Проверяем, что задачи были распределены
	if len(orch.tasks) != 0 {
		t.Errorf("Ожидалось, что все задачи будут распределены, осталось: %d", len(orch.tasks))
	}
}

func TestDistributeTasksNoAgents(t *testing.T) {
	// Создаем Orchestrator без агентов
	orch := &Orchestrator{
		tasks:   make(chan models.Task, 10),
		results: make(chan models.Result, 10),
		agents:  make(map[string]AgentInfo),
	}

	// Запускаем распределение задач в отдельной горутине
	go orch.distributeTasks()

	// Добавляем задачу в канал
	task := models.Task{ID: "task1", Numbers: []float64{1, 2}, Operators: []string{"+"}}
	orch.tasks <- task

	// Даем время для обработки задачи
	time.Sleep(1 * time.Second)

	// Проверяем, что задача не была распределена
	orch.mu.Lock()
	defer orch.mu.Unlock()

	if len(orch.tasks) != 0 {
		t.Errorf("Ожидалось, что задача будет удалена из канала, осталось: %d", len(orch.tasks))
	}
}

func TestDistributeTasksSingleAgent(t *testing.T) {
	// Создаем Orchestrator с одним агентом
	orch := &Orchestrator{
		tasks:   make(chan models.Task, 10),
		results: make(chan models.Result, 10),
		agents:  make(map[string]AgentInfo),
	}

	// Добавляем одного агента
	orch.mu.Lock()
	orch.agents["agent1"] = AgentInfo{ID: "agent1", ComputingPower: 10, LastSeen: time.Now()}
	orch.mu.Unlock()

	// Запускаем распределение задач в отдельной горутине
	go orch.distributeTasks()

	// Добавляем задачу в канал
	task := models.Task{ID: "task1", Numbers: []float64{1, 2}, Operators: []string{"+"}}
	orch.tasks <- task

	// Даем время для обработки задачи
	time.Sleep(1 * time.Second)

	// Проверяем, что задача была распределена
	orch.mu.Lock()
	defer orch.mu.Unlock()

	if len(orch.tasks) != 0 {
		t.Errorf("Ожидалось, что задача будет распределена, осталось: %d", len(orch.tasks))
	}
}
