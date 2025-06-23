package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akinoccc/hysaif/api/internal/context"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// Claims JWT声明
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthRequired 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authorization, "Bearer ") {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "无效的认证令牌格式"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		// 解析JWT
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "无效的认证令牌"})
			c.Abort()
			return
		}

		// 验证用户是否存在且状态正常
		var user models.User
		if err := models.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "用户不存在"})
			c.Abort()
			return
		}

		// 检查用户状态
		if !user.IsActive() {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "用户账号已被禁用或锁定"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user", user)

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return gin.Recovery()
}

// AuditLog 审计日志中间件
func AuditLog(action, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录审计日志
		user := context.GetCurrentUser(c)
		if user == nil {
			return
		}

		resourceID := c.Param("id")
		if resourceID == "" {
			resourceID = "N/A"
		}
	}
}

// RequirePermission 权限检查中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		user := context.GetCurrentUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "用户信息不存在"})
			c.Abort()
			return
		}

		// 检查权限
		if !user.HasPermission(resource, action) {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "权限不足，无法访问该资源"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(user models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
