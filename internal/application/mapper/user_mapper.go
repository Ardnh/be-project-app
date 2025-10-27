package mapper

import (
	"github.com/Ardnh/be-project-app/internal/application/dto"
	"github.com/Ardnh/be-project-app/internal/domain/entities"
)

func ToUserDTO(user *entities.User) *dto.UserDto {

	idString := (user.ID).String()
	return &dto.UserDto{
		ID:       idString,
		Username: user.Username,
		Email:    user.Email,
	}
}

func ToUserEntities(user *dto.CreateUserDto) *entities.User {
	return &entities.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}
