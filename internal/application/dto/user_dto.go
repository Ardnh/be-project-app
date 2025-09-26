package dto

type CreateUserDto struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=500"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
