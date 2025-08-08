package middleware

import (
	"fmt"
	"infotecstechtask/internal/models"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Миддлвар для валидации JSON объектов
// Внутри происходит сборка объекта модели и его валидация
// Если во время сборки или валидации возникает ошибка, конструируется ответ и отправляется клиенту
func JSONValidation(model any, validate *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		if c.ContentType() != "application/json" {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
				"error": "Only application/json content type is accepted for POST requests",
			})
			return
		}

		val := createModelInstance(model)

		if err := c.ShouldBindJSON(val); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				models.ValidationError{
					Error:   "Invalid JSON format",
					Details: formatJSONValidationErrors(err),
				},
			)
			return
		}

		if err := validate.Struct(val); err != nil {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				models.ValidationError{
					Error:   "Validation failed",
					Details: formatJSONValidationErrors(err),
				},
			)
			return
		}

		c.Set("validatedBody", val)
		c.Next()
	}
}

// Миддлвар для валидации Path и Query параметров
// Внутри происходит сборка объекта модели параметров запроса и его валидация
// Если во время сборки или валидации возникает ошибка, конструируется ответ и отправляется клиенту
func ParamsValidation(model any, validate *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		val := createModelInstance(model)

		if err := c.ShouldBindUri(val); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid path parameters",
				"details": formatErrors(err),
			})
			return
		}

		if err := c.ShouldBindQuery(val); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid query parameters",
			})
			return
		}

		if err := validate.Struct(val); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": formatJSONValidationErrors(err),
			})
			return
		}

		c.Set("validatedParams", val)
		c.Next()
	}
}

// Функция для форматирования ошибок валидации
func formatJSONValidationErrors(err error) []models.FieldError {
	errors := make([]models.FieldError, 0)
	for _, fieldErr := range err.(validator.ValidationErrors) {
		errors = append(
			errors,
			models.FieldError{
				Field:   fieldErr.Field(),
				Message: getValidationMessage(fieldErr),
			},
		)
	}

	return errors
}

// Функция для форматирования ошибок валидации query и path параметров запроса
func formatErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	}
	
	return errors
}

// Функция для создания инстанса модели
// Использует рефлексию
func createModelInstance(model any) any {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return reflect.New(t).Interface()
}

// Функция возвращает человекочитаемые ошибки валидации в зависимости от типа ошибки
func getValidationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "Field is required"
	case "uuid":
		return "Field must be a valid UUID"
	case "min":
		return fmt.Sprintf("Field must be greater than %s", fieldErr.Param())
	default:
		return fieldErr.Tag()
	}
}
