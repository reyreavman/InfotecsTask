package models

type Status string

// Возможные статусы транзакций
const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)