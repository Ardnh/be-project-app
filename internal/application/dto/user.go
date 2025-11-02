package dto

// ================ DTO ====================
type UserDto struct {
	ID       string `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=500"`
}

// ================ REQUEST DTO ====================
type CreateUserDto struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdateUserDto struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=500"`
}

// ================ RESPONSE DTO ====================

// ================ MAPPER ====================
func ToUserResponse() {

}
