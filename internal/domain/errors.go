package domain

import "errors"

var (
	// Auth & User
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrTokenExpired       = errors.New("token expired")

	// Validation
	ErrInvalidInput = errors.New("invalid input")

	// Domain
	ErrProjectNotFound = errors.New("project not found")
	ErrDuplicateEntry  = errors.New("duplicate entry")

	// Generic
	ErrInternalServerError = errors.New("internal server error")
)
