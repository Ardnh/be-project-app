package users

import (
	"errors"
	"fmt"
	"project-app/helper"
	"project-app/model"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	Register(ctx *fiber.Ctx, req *model.User) (*string, error)
	FindByUsernameOrEmail(ctx *fiber.Ctx, req string, isEmail bool) (*model.User, error)
	FindByEmail(ctx *fiber.Ctx, email string) (*model.User, error)
	CreatUserProfileById(ctx *fiber.Ctx, req *model.ProfileCreateRequest) error
	UpdateProfileById(ctx *fiber.Ctx, userId uuid.UUID, req model.ProfileUpdateRequest) error
	GetProfileById(ctx *fiber.Ctx, userId uint) (*model.Profile, error)
}

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {

	return &UsersRepositoryImpl{
		Db: db,
	}
}

var tableUser = "users"
var tableProfile = "user_profiles"
var tableFollowers = "follow_users"

func (repository *UsersRepositoryImpl) GetProfileById(ctx *fiber.Ctx, userId uint) (*model.Profile, error) {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	var result model.Profile
	err := tx.
		WithContext(ctx.Context()).
		Table(tableProfile).
		Where("user_id = ?", userId).
		Take(&result)

	if err.Error != nil {
		return nil, err.Error
	}

	return &result, nil
}

func (repository *UsersRepositoryImpl) Register(ctx *fiber.Ctx, req *model.User) (*string, error) {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// 1. Insert user to user table
	result := tx.
		WithContext(ctx.Context()).
		Table(tableUser).
		Create(&req)

	if result.Error != nil {
		return nil, result.Error
	}

	return &req.ID, nil
}

func (repository *UsersRepositoryImpl) FindByUsernameOrEmail(ctx *fiber.Ctx, req string, isEmail bool) (*model.User, error) {
	var result model.User

	// Buat query dasar
	query := repository.Db.WithContext(ctx.Context()).Table(tableUser)

	// Tentukan kondisi berdasarkan isEmail
	if isEmail {
		query = query.Where("email = ?", req)
	} else {
		query = query.Where("username = ?", req)
	}

	// Jalankan query untuk mendapatkan hasil pertama
	err := query.First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repository *UsersRepositoryImpl) FindByEmail(ctx *fiber.Ctx, email string) (*model.User, error) {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	var result model.User

	query := tx.WithContext(ctx.Context()).Table(tableUser)
	query = query.Where("email = ?", strings.ToLower(email)).Find(&result)

	err := query.Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &result, nil
}

func (repository *UsersRepositoryImpl) FindFollowersByUserId(ctx *fiber.Ctx, userId uint, page int, pageSize int, searchQuery string) ([]model.UserWithProfile, int64, error) {

	var userWithProfile []model.UserWithProfile
	var totalCount int64

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	// Offset
	offset := (page - 1) * pageSize

	// Query
	query := tx.WithContext(ctx.Context()).Table(tableFollowers)
	if searchQuery != "" {
		query = query.Where("username LIKE ? ", "%"+searchQuery+"%")
	}

	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	errResult := query.
		Offset(offset).
		Limit(pageSize).
		Find(&userWithProfile).
		Error

	if errResult != nil {
		return nil, 0, errResult
	}

	return userWithProfile, totalCount, nil
}

func (repository *UsersRepositoryImpl) UpdateProfileById(ctx *fiber.Ctx, userId uuid.UUID, req model.ProfileUpdateRequest) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	fmt.Println(userId)

	err := tx.WithContext(ctx.Context()).Table(tableProfile).Where("user_id = ?", userId).Updates(&req)

	if err != nil {
		return err.Error
	}

	return nil
}

func (repository *UsersRepositoryImpl) CreatUserProfileById(ctx *fiber.Ctx, req *model.ProfileCreateRequest) error {

	tx := repository.Db.Begin()
	defer helper.CommitOrRollback(tx)

	err := tx.WithContext(ctx.Context()).Table(tableProfile).Create(&req).Error

	if err != nil {
		return err
	}

	return nil
}

// func (repository *UsersRepositoryImpl) DeleteFollowUserById(followId int) error {

// }
