package users

import (
	"fmt"
	"project-app/helper"
	userRepository "project-app/repository/users"
	"strconv"

	"project-app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UsersHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	FindUserProfileById(c *fiber.Ctx) error
	UpdateProfileById(c *fiber.Ctx) error
	GetProfileById(c *fiber.Ctx) error
}

type UsersHandlerImpl struct {
	UsersRepository userRepository.UsersRepository
	Validate        *validator.Validate
}

func NewUsersHandler(db *gorm.DB, validate *validator.Validate) UsersHandler {
	user := userRepository.NewUsersRepository(db)
	return &UsersHandlerImpl{
		UsersRepository: user,
		Validate:        validate,
	}
}

// Login user
// @Summary Login user
// @Description Login user
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "Login"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/login [post]
func (handler *UsersHandlerImpl) Login(c *fiber.Ctx) error {

	// Read body request
	var request model.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// Validate incoming request
	errValidate := handler.Validate.Struct(request)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// Get user by email
	userResult, errUserFindByEmail := handler.UsersRepository.FindByEmail(c, request.Email)
	if errUserFindByEmail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errUserFindByEmail.Error(),
		})
	}

	fmt.Println(userResult)

	// compare password from body request and from database
	errComparePassword := helper.CheckPasswordHash(request.Password, userResult.Password)
	if !errComparePassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Wrong password!",
		})
	}

	fmt.Println("login result")
	fmt.Println(request.Password)
	fmt.Println("err compare")
	fmt.Println(errComparePassword)

	// generate jwt token
	token, errGenerateToken := helper.GenerateToken(userResult.ID)
	if errGenerateToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errGenerateToken.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Login successfully",
		"data": fiber.Map{
			"token":    token,
			"username": userResult.Username,
		},
	})
}

// Register user
// @Summary Register user
// @Description Register user
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "Login"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/register [post]
func (handler *UsersHandlerImpl) Register(c *fiber.Ctx) error {

	// 1. Parser body request
	var request model.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	// 2. Validasi json yang dikirim
	errValidate := handler.Validate.Struct(request)
	if errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": errValidate.Error(),
		})
	}

	// 3. Cek apakah user dengan email yang dikirim sudah ada di database
	result, err := handler.UsersRepository.FindByEmail(c, request.Email)
	if result.Email == request.Email {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "User already exist",
		})
	}

	// 4. Hash user password
	hashResult, errHash := helper.HashPassword(request.Password)
	if errHash != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Error hashing password",
		})
	}
	request.Password = hashResult

	// 5. Save user to database
	req := model.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	userId, errRegister := handler.UsersRepository.Register(c, &req)
	if errRegister != nil && userId == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	defaultDataToProfile := model.ProfileCreateRequest{
		UserId:    req.ID,
		Bio:       "",
		Role:      "",
		Facebook:  "",
		Instagram: "",
		Linkedin:  "",
		Twitter:   "",
	}

	errProfile := handler.UsersRepository.CreatUserProfileById(c, &defaultDataToProfile)

	if errProfile != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errProfile.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully register user",
	})
}

// Find user profile by id
// @Summary Find user profile by id
// @Description Find user profile by id
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/following [get]
func (handler *UsersHandlerImpl) FindUserProfileById(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Hello from",
	})
}

// Update profile by id
// @Summary Update profile by id
// @Description Update profile by id
// @Security Bearer
// @Tags Users
// @Accept json
// @Produce json
// @Param body body model.ProfileUpdateRequestBody true "Update profile"
// @Success 200 {object} map[string]interface{} "Success update category"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [put]
func (handler *UsersHandlerImpl) UpdateProfileById(c *fiber.Ctx) error {

	var request model.ProfileUpdateRequestBody
	userId := helper.UserId

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	reqData := model.ProfileUpdateRequest{
		Bio:       request.Bio,
		Role:      request.Role,
		Facebook:  request.Facebook,
		Instagram: request.Instagram,
		Linkedin:  request.LinkedIn,
		Twitter:   request.Twitter,
	}

	err := handler.UsersRepository.UpdateProfileById(c, userId, reqData)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully update user profile",
	})
}

// Get profile by id
// @Summary Get profile by id
// @Description Get profile by id
// @Tags Users
// @Produce json
// @Security Bearer
// @Param userId path string true "userId"
// @Success 200 {object} map[string]interface{} "Success get profile by id"
// @Failure 400 {object} map[string]interface{} "Invalid request body or missing required fields"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile/:userId [get]
func (handler *UsersHandlerImpl) GetProfileById(c *fiber.Ctx) error {
	userId := c.Params("userId", "")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid user id",
		})
	}

	userIdVal, err := strconv.ParseUint(userId, 10, 32) // basis 10, 32-bit
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid user id",
		})
	}

	result, err := handler.UsersRepository.GetProfileById(c, uint(userIdVal))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Successfully get profile",
		"data":    result,
	})
}
