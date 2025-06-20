package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/llamacto/llama-gin-kit/pkg/logger"
)

// UserHandler 用户处理器
type UserHandler struct {
	service *UserServiceImpl
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(service *UserServiceImpl) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserRegisterRequest true "注册信息"
// @Success 200 {object} User
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserLoginRequest true "登录信息"
// @Success 200 {object} User
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户的个人资料
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserUpdateRequest true "用户信息"
// @Success 200 {object} User
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	userID := userIDVal.(uint)

	var req UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	user, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserChangePasswordRequest true "密码信息"
// @Success 200 {string} string "密码修改成功"
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	userID := userIDVal.(uint)

	var req UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 通过邮箱重置用户密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body UserPasswordResetRequest true "邮箱信息"
// @Success 200 {string} string "重置密码邮件已发送"
// @Router /users/password/reset [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req UserPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置邮件已发送"})
}

// GetProfile 获取用户个人资料
// @Summary 获取用户个人资料
// @Description 获取当前登录用户的个人资料
// @Tags 用户
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} User
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	user, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteAccount 删除账户
// @Summary 删除账户
// @Description 删除当前用户的账户
// @Tags 用户
// @Success 200 {string} string "账户已删除"
// @Router /users/account [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}
	userID := userIDVal.(uint)

	if err := h.service.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "账户已删除"})
}

// Get 获取指定用户信息
// @Summary 获取指定用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} User
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// List 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {array} User
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		logger.Error("获取用户列表失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total, "list": users})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserInfo
// @Router /users/info/{id} [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userInfo, err := h.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
