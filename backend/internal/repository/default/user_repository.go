package repository

import (
	"context"

	dto "github.com/hadi-projects/go-react-starter/internal/dto/default"
	entity "github.com/hadi-projects/go-react-starter/internal/entity/default"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.User, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByEmailSimple(ctx context.Context, email string) (*entity.User, error)
	FindRoleByName(ctx context.Context, name string) (*entity.Role, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindAll(ctx context.Context, pagination *dto.PaginationRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{})

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query = query.Where("email LIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (pagination.GetPage() - 1) * pagination.GetLimit()
	err := query.Order("id DESC").
		Preload("Role").
		Limit(pagination.GetLimit()).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Role.Permissions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Preload("Role.Permissions").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmailSimple(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}
