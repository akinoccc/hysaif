package validation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/akinoccc/hysaif/api/types"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// HandleValidationErrors 处理gin的validation错误并返回详细的错误信息
func HandleValidationErrors(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make(map[string][]string)

		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			message := GetValidationErrorMessage(fieldError)
			details[field] = append(details[field], message)
		}

		c.JSON(http.StatusBadRequest, types.ValidationErrorResponse{
			Error:   "请求参数验证失败",
			Details: details,
		})
		return
	}

	// 如果不是validation错误，返回通用错误
	c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误: " + err.Error()})
}

// GetValidationErrorMessage 根据validation错误类型返回中文错误信息
func GetValidationErrorMessage(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s是必填字段", field)
	case "email":
		return fmt.Sprintf("%s必须是有效的邮箱地址", field)
	case "min":
		return fmt.Sprintf("%s长度不能少于%s个字符", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s长度不能超过%s个字符", field, fe.Param())
	case "oneof":
		values := strings.Split(fe.Param(), " ")
		return fmt.Sprintf("%s必须是以下值之一: %s", field, strings.Join(values, ", "))
	case "dive":
		return fmt.Sprintf("%s包含无效的元素", field)
	default:
		return fmt.Sprintf("%s验证失败", field)
	}
}
