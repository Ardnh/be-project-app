package model

import "gorm.io/gorm"

type Category struct {
	*gorm.Model
	Name     string `gorm:"type:varchar(100)"`
	Projects []Projects
}

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryDeleteRequest struct {
	Id int `json:"id" validate:"required,number"`
}
