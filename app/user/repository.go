package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	Get(ctx context.Context, id uint) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]*User, int64, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByID(id string) (*UserInfo, error)
}

// UserRepositoryImpl 用户仓库实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

// Get 获取用户
func (r *UserRepositoryImpl) Get(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
	var users []*User
	var total int64

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(id string) (*UserInfo, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        fmt.Sprintf("%d", user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
		Bio:       user.Bio,
		Status:    user.Status,
		LastLogin: user.LastLogin,
	}, nil
}
