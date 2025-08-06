package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid JSON format",
				"details": formatValidationErrors(err),
			})
			return
		}

		if err := validate.Struct(val); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "Validation failed",
				"details": formatValidationErrors(err),
			})
			return
		}

		c.Set("validatedBody", val)
		c.Next()
	}
}

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
				"details": formatValidationErrors(err),
			})
			return
		}
		
		// val := &models.GetWalletBalanceRequest{
		// 	ID: uuid.MustParse("b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a10"),
		// }

		if err := c.ShouldBindQuery(val); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid query parameters",
				"details": formatValidationErrors(err),
			})
			return
		}		

		if err := validate.Struct(val); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "Validation failed",
				"details": formatValidationErrors(err),
			})
			return
		}

		c.Set("validatedParams", val)
		c.Next()
	}
}

func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	} else {
		errors["_error"] = err.Error()
	}
	return errors
}

func createModelInstance(model any) any {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return reflect.New(t).Interface()
}
