package repository

import (
	"errors"

	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Struct implementasi (private, lowercase)
type userRepositoryImpl struct {
	db    *gorm.DB
	redis *redis.Client
}

// Constructor yang return INTERFACE dari domain
// ⬇️ Return type adalah INTERFACE dari domain
func NewUserRepository(db *gorm.DB, redis *redis.Client) repositories.UserRepository {
	return &userRepositoryImpl{
		db:    db,
		redis: redis,
	}
}

func (r *userRepositoryImpl) FindByID(id string) (*entities.User, error) {

	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindAll(limit int, offset int) ([]*entities.User, error) {

	var users []*entities.User
	query := r.db.Model(&entities.User{})

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) Create(user *entities.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
}

func (r *userRepositoryImpl) Update(user *entities.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(user).Error
	})
}

func (r *userRepositoryImpl) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := r.db.Where("id = ?", id).Delete(&entities.User{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found")
		}
		return nil
	})
}

func (r *userRepositoryImpl) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
