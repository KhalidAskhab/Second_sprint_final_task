package models

type Task struct {
	ID        string    `json:"id"`
	Numbers   []float64 `json:"numbers"`
	Operators []string  `json:"operators"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type Expression struct {
	ID         string  `json:"id"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result,omitempty"`
}
