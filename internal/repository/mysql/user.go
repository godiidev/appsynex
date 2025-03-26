package mysql

import (
	"github.com/godiidev/appsynex/internal/domain/models"
	"github.com/godiidev/appsynex/internal/repository/interfaces"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(page, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	query := r.db.Model(&models.User{})

	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Paginate
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("Roles").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) FindByIDWithRoles(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// Begin transaction
	tx := r.db.Begin()

	// Remove existing roles
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new roles
	for _, roleID := range roleIDs {
		if err := tx.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}
