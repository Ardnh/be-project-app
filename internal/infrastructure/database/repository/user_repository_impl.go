package repository

import (
	"context"
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

func (r *userRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.User, error) {

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

func (r *userRepositoryImpl) FindAll(ctx context.Context, limit int, offset int) ([]*entities.User, error) {

	var users []*entities.User
	query := r.db.Model(&entities.User{})

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(user).Error
	})
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id string) error {
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

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user *entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
