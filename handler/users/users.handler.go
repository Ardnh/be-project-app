package users

import (
	"project-app/helper"
	userRepository "project-app/repository/users"
	"strconv"

	"project-app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

func NewUsersHandler(userRepository userRepository.UsersRepository, validate *validator.Validate) UsersHandler {
	return &UsersHandlerImpl{
		UsersRepository: userRepository,
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

	start := helper.StartTime()
	// Logging start of login process
	logrus.Info("Login request received")

	// Read body request
	var request model.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse login request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     err.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.WithFields(logrus.Fields{
		"email": request.Email,
	}).Info("Request body parsed successfully")

	// Validate incoming request
	errValidate := handler.Validate.Struct(request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in login request")

		return c.Status(fiber.StatusBadRequest).JSON(helper.HandleValidationError(errValidate, start, request))
	}

	email := model.IsEmail{
		Email: request.Email,
	}

	var isEmail bool = false
	errIsEmail := handler.Validate.Struct(email)
	if errIsEmail != nil {
		logrus.WithFields(logrus.Fields{
			"email": request.Email,
		}).Warn("Invalid email format")
		isEmail = false
	} else {
		isEmail = true
	}

	logrus.WithFields(logrus.Fields{
		"email":    request.Email,
		"is_email": isEmail,
	}).Info("Email validation complete")

	// Get user by email
	userResult, errUserFindByEmail := handler.UsersRepository.FindByUsernameOrEmail(c, request.Email, isEmail)
	if errUserFindByEmail != nil {
		logrus.WithError(errUserFindByEmail).Error("Failed to find user by email or username")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     errUserFindByEmail.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  userResult.ID,
		"username": userResult.Username,
	}).Info("User found successfully")

	// compare password from body request and from database
	errComparePassword := helper.CheckPasswordHash(request.Password, userResult.Password)
	if !errComparePassword {
		logrus.Warn("Password mismatch for user", userResult.Username)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":        fiber.StatusBadRequest,
			"message":     "Wrong password!",
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.Info("Password matched, generating JWT token")

	// generate jwt token
	token, errGenerateToken := helper.GenerateToken(userResult.ID)
	if errGenerateToken != nil {
		logrus.WithError(errGenerateToken).Error("Failed to generate JWT token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     errGenerateToken.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.WithFields(logrus.Fields{
		"user_id": userResult.ID,
		"token":   token,
	}).Info("JWT token generated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        fiber.StatusOK,
		"message":     "Login successfully",
		"processTime": helper.ElapsedTime(start),
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

	start := helper.StartTime()
	// Logging start of registration process
	logrus.Info("Register request received")

	// 1. Parse body request
	var request model.RegisterRequest
	if err := c.BodyParser(&request); err != nil {
		logrus.WithError(err).Error("Failed to parse register request body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     err.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.WithFields(logrus.Fields{
		"email":    request.Email,
		"username": request.Username,
	}).Info("Request body parsed successfully")

	// 2. Validate json yang dikirim
	errValidate := handler.Validate.Struct(request)
	if errValidate != nil {
		logrus.WithError(errValidate).Warn("Validation error in register request")
		return c.Status(fiber.StatusBadRequest).JSON(helper.HandleValidationError(errValidate, start, request))
	}

	// 3. Cek apakah user dengan email yang dikirim sudah ada di database
	result, err := handler.UsersRepository.FindByEmail(c, request.Email)
	if err != nil {
		logrus.WithError(err).Error("Failed to check existing user by email")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     err.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	if result.Email == request.Email {
		logrus.WithFields(logrus.Fields{
			"email": request.Email,
		}).Warn("User already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":        fiber.StatusBadRequest,
			"message":     "User already exist",
			"processTime": helper.ElapsedTime(start),
		})
	}

	// 4. Hash user password
	hashResult, errHash := helper.HashPassword(request.Password)
	if errHash != nil {
		logrus.WithError(errHash).Error("Error hashing password")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     "Error hashing password",
			"processTime": helper.ElapsedTime(start),
		})
	}
	request.Password = hashResult

	logrus.Info("Password hashed successfully")

	// 5. Save user to database
	req := model.User{
		ID:       uuid.NewString(),
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	userId, errRegister := handler.UsersRepository.Register(c, &req)
	if errRegister != nil && userId == nil {
		logrus.WithError(errRegister).Error("Failed to register user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":        fiber.StatusInternalServerError,
			"message":     errRegister.Error(),
			"processTime": helper.ElapsedTime(start),
		})
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  req.ID,
		"username": req.Username,
		"email":    req.Email,
	}).Info("User registered successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":        fiber.StatusOK,
		"message":     "Successfully registered user",
		"processTime": helper.ElapsedTime(start),
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
