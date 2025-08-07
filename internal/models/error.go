package models

type Error struct {
	Error string
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Error   string
	Details []FieldError
}
