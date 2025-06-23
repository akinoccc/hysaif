package handlers

import (
	"net/http"
	"time"

	"github.com/akinoccc/hysaif/api/middleware"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/types"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login 用户登录
func Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ? AND status = ?", req.Email, "active").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "生成令牌失败"})
		return
	}

	// 记录登录时间和IP地址
	user.LastLoginAt = uint64(time.Now().UnixMilli())
	user.LastLoginIP = c.ClientIP()
	models.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": user.LastLoginAt,
		"last_login_ip": user.LastLoginIP,
	})

	c.JSON(http.StatusOK, types.LoginResponse{
		Token: token,
		User:  user,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, types.MessageResponse{Message: "登出成功"})
}
