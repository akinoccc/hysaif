package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/akinoccc/hysaif/api/middleware"
	"github.com/akinoccc/hysaif/api/models"
	webauthnPkg "github.com/akinoccc/hysaif/api/packages/webauthn"
	"github.com/akinoccc/hysaif/api/types"

	pkgContext "github.com/akinoccc/hysaif/api/packages/context"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/bcrypt"
)

// 存储会话数据的临时存储（实际生产中应使用 Redis 等）
var (
	registrationSessions = make(map[string]*webauthn.SessionData)
	loginSessions        = make(map[string]*webauthn.SessionData)
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
	if err := models.DB.Where("email = ? AND status = ?", req.Email, models.StatusActive).First(&user).Error; err != nil {
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

// WebAuthnBeginRegistration 开始 WebAuthn 注册
func WebAuthnBeginRegistration(c *gin.Context) {
	var req struct {
		CredentialName string `json:"credential_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 获取当前用户
	user := pkgContext.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取 WebAuthn 实例
	webAuthn, err := webauthnPkg.GetWebAuthn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "WebAuthn 初始化失败"})
		return
	}

	// 开始注册流程
	options, session, err := webAuthn.BeginRegistration(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "开始注册失败"})
		return
	}

	// 保存会话数据
	registrationSessions[user.ID] = session

	c.JSON(http.StatusOK, options)
}

// WebAuthnFinishRegistration 完成 WebAuthn 注册
func WebAuthnFinishRegistration(c *gin.Context) {
	// 获取当前用户
	user := pkgContext.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取会话数据
	session, ok := registrationSessions[user.ID]
	if !ok {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "无效的会话"})
		return
	}

	// 先读取原始请求体
	bodyBytes, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "读取请求数据失败"})
		return
	}

	// 解析请求体以获取credential name
	var reqBody struct {
		CredentialName string `json:"credential_name"`
	}
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	if reqBody.CredentialName == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "凭证名称不能为空"})
		return
	}

	// 重新创建请求对象给FinishRegistration使用
	newReq := c.Request.Clone(c.Request.Context())
	newReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 获取 WebAuthn 实例
	webAuthn, err := webauthnPkg.GetWebAuthn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "WebAuthn 初始化失败"})
		return
	}

	// 解析凭证
	credential, err := webAuthn.FinishRegistration(user, *session, newReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "注册失败: " + err.Error()})
		return
	}

	// 保存凭证
	if err := user.AddWebAuthnCredential(credential, reqBody.CredentialName); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "保存凭证失败"})
		return
	}

	// 清理会话
	delete(registrationSessions, user.ID)

	c.JSON(http.StatusOK, types.MessageResponse{Message: "注册成功"})
}

// WebAuthnBeginLogin 开始 WebAuthn 登录
func WebAuthnBeginLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ? AND status = ?", req.Email, models.StatusActive).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "用户不存在"})
		return
	}

	// 获取 WebAuthn 实例
	webAuthn, err := webauthnPkg.GetWebAuthn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "WebAuthn 初始化失败"})
		return
	}

	// 开始登录流程
	options, session, err := webAuthn.BeginLogin(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "开始登录失败"})
		return
	}

	// 保存会话数据
	loginSessions[user.ID] = session

	c.JSON(http.StatusOK, options)
}

// WebAuthnFinishLogin 完成 WebAuthn 登录
func WebAuthnFinishLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "请求参数错误"})
		return
	}

	log.Printf("WebAuthn login request: %+v", req)

	// 查找用户
	var user models.User
	if err := models.DB.Where("email = ? AND status = ?", req.Email, models.StatusActive).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "用户不存在"})
		return
	}

	// 获取会话数据
	session, ok := loginSessions[user.ID]
	if !ok {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "无效的会话"})
		return
	}

	// 获取 WebAuthn 实例
	webAuthn, err := webauthnPkg.GetWebAuthn()
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "WebAuthn 初始化失败"})
		return
	}

	// 完成登录
	credential, err := webAuthn.FinishLogin(&user, *session, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "登录失败: " + err.Error()})
		return
	}

	// 更新凭证信息（签名计数等）
	if err := user.UpdateWebAuthnCredential(credential); err != nil {
		// 记录错误但不影响登录
		// log.Printf("更新凭证失败: %v", err)
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

	// 清理会话
	delete(loginSessions, user.ID)

	c.JSON(http.StatusOK, types.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetUserCredentials 获取用户的 WebAuthn 凭证列表
func GetUserCredentials(c *gin.Context) {
	// 获取当前用户
	user := pkgContext.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 获取用户的凭证
	var credentials []models.WebAuthnCredential
	if err := models.DB.Where("user_id = ?", user.ID).Find(&credentials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "获取凭证失败"})
		return
	}

	c.JSON(http.StatusOK, credentials)
}

// DeleteUserCredential 删除用户的 WebAuthn 凭证
func DeleteUserCredential(c *gin.Context) {
	credentialID := c.Param("id")

	// 获取当前用户
	user := pkgContext.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "未授权"})
		return
	}

	// 删除凭证
	result := models.DB.Where("id = ? AND user_id = ?", credentialID, user.ID).Delete(&models.WebAuthnCredential{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "删除凭证失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "凭证不存在"})
		return
	}

	c.JSON(http.StatusOK, types.MessageResponse{Message: "删除成功"})
}
