package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/llamacto/llama-gin-kit/pkg/email"
	"github.com/llamacto/llama-gin-kit/pkg/jwt"
	"github.com/llamacto/llama-gin-kit/pkg/logger"
	"github.com/llamacto/llama-gin-kit/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService User 服务接口
type UserService interface {
	Create(ctx context.Context, model *User) error
	Update(ctx context.Context, model *User) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]*User, int64, error)
	Register(req *UserRegisterRequest) (*User, error)
	Login(req *UserLoginRequest) (*UserLoginResponse, error)
	UpdateProfile(userID uint, req *UserUpdateRequest) (*User, error)
	ChangePassword(userID uint, req *UserChangePasswordRequest) error
	ResetPassword(req *UserPasswordResetRequest) error
	GetProfile(userID uint) (*User, error)
	DeleteAccount(userID uint) error
	GetUserByID(id uint) (*UserInfo, error)
	GetByID(id uint) (*User, error)
}

// UserServiceImpl User 服务实现
type UserServiceImpl struct {
	repo UserRepository
}

// NewUserService 创建 User 服务
func NewUserService(repo UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

// Create 创建 User
func (s *UserServiceImpl) Create(ctx context.Context, model *User) error {
	return s.repo.Create(ctx, model)
}

// Update 更新 User
func (s *UserServiceImpl) Update(ctx context.Context, model *User) error {
	return s.repo.Update(ctx, model)
}

// Delete 删除 User
func (s *UserServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// Get 获取 User
func (s *UserServiceImpl) Get(ctx context.Context, id uint) (*User, error) {
	return s.repo.Get(ctx, id)
}

// List 获取 User 列表
func (s *UserServiceImpl) List(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
	return s.repo.List(ctx, page, pageSize)
}

// Register 用户注册
func (s *UserServiceImpl) Register(req *UserRegisterRequest) (*User, error) {
	ctx := context.Background()

	// 检查邮箱是否已存在
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 发送欢迎邮件
	if err := email.SendWelcomeEmail(user.Email, user.Username); err != nil {
		logger.Error("发送欢迎邮件失败:", err)
	}

	return user, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(req *UserLoginRequest) (*UserLoginResponse, error) {
	ctx := context.Background()

	// Try to find user by username first
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		// If not found by username, try email
		user, err = s.repo.GetByEmail(ctx, req.Username)
		if err != nil {
			return nil, errors.New("用户名或密码错误")
		}
	}

	if user.Status == 0 {
		return nil, errors.New("账户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}

	now := time.Now()
	user.LastLogin = &now
	if err := s.repo.Update(ctx, user); err != nil {
		logger.Error("更新用户最后登录时间失败:", err)
	}

	return &UserLoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// UpdateProfile 更新用户信息
func (s *UserServiceImpl) UpdateProfile(userID uint, req *UserUpdateRequest) (*User, error) {
	ctx := context.Background()

	user, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(userID uint, req *UserChangePasswordRequest) error {
	ctx := context.Background()

	user, err := s.repo.Get(ctx, userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// ResetPassword 重置密码
func (s *UserServiceImpl) ResetPassword(req *UserPasswordResetRequest) error {
	ctx := context.Background()

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("邮箱不存在")
	}

	// 生成随机密码
	newPassword := utils.GenerateRandomString(12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	// 发送重置密码邮件
	if err := email.SendPasswordResetEmail(user.Email, newPassword); err != nil {
		logger.Error("发送重置密码邮件失败:", err)
		return errors.New("发送重置密码邮件失败")
	}

	return nil
}

// GetProfile 获取用户信息
func (s *UserServiceImpl) GetProfile(userID uint) (*User, error) {
	ctx := context.Background()
	user, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// DeleteAccount 删除账户
func (s *UserServiceImpl) DeleteAccount(userID uint) error {
	ctx := context.Background()
	if err := s.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("删除账户失败: %w", err)
	}
	return nil
}

// GetUserByID retrieves user information by ID.
func (s *UserServiceImpl) GetUserByID(id uint) (*UserInfo, error) {
	return s.repo.FindByID(id)
}

// GetByID retrieves a user by their ID.
func (s *UserServiceImpl) GetByID(id uint) (*User, error) {
	ctx := context.Background()
	return s.repo.Get(ctx, id)
}
