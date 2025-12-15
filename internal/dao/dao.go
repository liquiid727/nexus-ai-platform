package dao

import (
	"context"
	"next-ai-gateway/internal/model/entity"

	"gorm.io/gorm"
)

type DepartmentDAO struct {
	db *gorm.DB
}

func NewDepartmentDAO(db *gorm.DB) *DepartmentDAO {
	return &DepartmentDAO{db: db}
}

func (d *DepartmentDAO) Create(ctx context.Context, dept *entity.Department) error {
	return d.db.WithContext(ctx).Create(dept).Error
}

func (d *DepartmentDAO) GetByID(ctx context.Context, id string) (*entity.Department, error) {
	var dept entity.Department
	if err := d.db.WithContext(ctx).Where("id = ?", id).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (d *DepartmentDAO) Update(ctx context.Context, dept *entity.Department) error {
	return d.db.WithContext(ctx).Save(dept).Error
}

func (d *DepartmentDAO) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&entity.Department{}, "id = ?", id).Error
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) Create(ctx context.Context, user *entity.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *UserDAO) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := d.db.WithContext(ctx).Where("email = ?", email).Preload("Department").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type AiModelDAO struct {
	db *gorm.DB
}

func NewAiModelDAO(db *gorm.DB) *AiModelDAO {
	return &AiModelDAO{db: db}
}

func (d *AiModelDAO) Create(ctx context.Context, model *entity.AiModel) error {
	return d.db.WithContext(ctx).Create(model).Error
}

func (d *AiModelDAO) GetAllActive(ctx context.Context) ([]entity.AiModel, error) {
	var models []entity.AiModel
	if err := d.db.WithContext(ctx).Where("is_active = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}
