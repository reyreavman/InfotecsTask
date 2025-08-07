package models

// Модель ошибки
type Error struct {
	Error string
}

// Модель ошибки валидации в конкретном поле
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Общая модель, состоящая из общего сообщения об ошибке и деталей
type ValidationError struct {
	Error   string
	Details []FieldError
}
