# Финальная задача спринта 2. Конкурентное программирование
## Компоненты системы

- **Оркестратор**: Это серверная часть, которая принимает арифметические выражения, разбивает их на задачи и отправляет их агентам для выполнения.
- **Агент**: Это демон, который получает задачи от оркестратора, выполняет их и возвращает результаты обратно оркестратору.

## Пакеты
- │Second_sprint_final_task/
- ├── cmd/
- │   ├── main.go
- │   └── agent/
- │       └── main.go
- ├── internal/
- │   ├── agent/
- │   │   ├── agent.go
- │   │   └── agent_test.go
- │   ├── orchestrator/
- │   │   ├── orchestrator.go
- │   │   └── orchestrator_test.go
- │   └── application/
- │       ├── application.go
- │       └── application_test.go
- ├── pkg/
- │   ├── calculation/
- │   │   ├── calculation.go
- │   │   ├── calculation_test.go
- │   │   └── errors.go
- │   └── models/
- │       └── models.go
- ├── go.mod
- ├── go.sum
- └── README.md

## Описание файлов и пакетов
-cmd/

-cmd/main.go: Основной файл для запуска приложения (например, сервера или оркестратора).

-cmd/agent/main.go: Основной файл для запуска агента.

-internal/

-internal/agent/agent.go: Логика агента, который выполняет задачи.

-internal/agent/agent_test.go: Тесты для агента.

-internal/orchestrator/orchestrator.go: Логика оркестратора, который распределяет задачи между агентами.

-internal/orchestrator/orchestrator_test.go: Тесты для оркестратора.

-internal/application/application.go: Логика приложения (например, HTTP-сервер или CLI).

-internal/application/application_test.go: Тесты для приложения.

-pkg/

-pkg/calculation/calculation.go: Логика для вычисления математических выражений.

-pkg/calculation/calculation_test.go: Тесты для модуля вычислений.

-pkg/calculation/errors.go: Определение ошибок для модуля вычислений.

-pkg/models/models.go: Модели данных, используемые в проекте (например, Task, Result, Expression)
