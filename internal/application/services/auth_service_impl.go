package services

import (
	"context"
	"errors"
	"time"

	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/config"
	"github.com/Ardnh/be-project-app/internal/domain"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
	"github.com/Ardnh/be-project-app/internal/domain/repositories"
	"github.com/Ardnh/be-project-app/internal/domain/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) services.AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *dto.LoginDto) (*string, error) {

	emailIsExist, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !emailIsExist {
		return nil, domain.ErrUserNotFound
	}

	user, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	// Generate jwt token
	// Load config
	config := config.LoadConfig()
	if config.App.JWTSecret == "" {
		return nil, errors.New("Failed to load jwt secret")
	}

	// Secret key untuk signing
	secretKey := []byte(config.App.JWTSecret)

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims (data dalam token)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expire 24 jam

	// Generate signed token string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *dto.RegisterDto) error {

	emailIsExist, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if emailIsExist {
		return domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userPassword := string(hashedPassword)
	userEntities := &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: userPassword,
	}

	errCreate := s.userRepo.Create(ctx, userEntities)
	if errCreate != nil {
		return errCreate
	}

	return nil
}
