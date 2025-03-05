package orchestrator

import (
	"SSFT/pkg/models"
	"log"
	"sync"
	"time"
)

type AgentInfo struct {
	ID             string
	ComputingPower int
	LastSeen       time.Time
}

type Orchestrator struct {
	tasks   chan models.Task
	results chan models.Result
	agents  map[string]AgentInfo
	mu      sync.Mutex
}

func (o *Orchestrator) distributeTasks() {
	for task := range o.tasks {
		o.mu.Lock()
		var bestAgent string
		maxPower := 0
		for id, agent := range o.agents {
			if agent.ComputingPower > maxPower {
				bestAgent = id
				maxPower = agent.ComputingPower
			}
		}
		o.mu.Unlock()

		if bestAgent != "" {
			log.Printf("Задача с ID %s распределена агенту %s", task.ID, bestAgent)
		} else {
			log.Printf("Нет доступных агентов для задачи с ID %s", task.ID)
		}
	}
}
