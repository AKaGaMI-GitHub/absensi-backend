package model

import (
	"time"
)

type UserRole string

type User struct {
	ID           string    `json:"id" bson:"_id"`
	NamaDepan    string    `json:"namaDepan" bson:"namaDepan"`
	NamaBelakang string    `json:"namaBelakang" bson:"namaBelakang"`
	Username     string    `json:"username" bson:"username"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"-" bson:"password"` // jangan expose ke JSON
	Avatar       string    `json:"avatar" bson:"avatar"`
	Role         string    `json:"-" bson:"role"`
	RoleUser     *RoleUser `json:"role,omitempty"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updated_at"`
}

type CreateUserInput struct {
	NamaDepan    string `json:"namaDepan" form:"namaDepan" validate:"required,min=2,max=40"`
	NamaBelakang string `json:"namaBelakang" form:"namaBelakang" validate:"max=60"`
	Username     string `json:"username" form:"username" validate:"required,min=5,max=30"`
	Email        string `json:"email" form:"email" validate:"required,email"`
	Password     string `json:"password" form:"password" validate:"required,min=6"`
	Role         string `json:"role" form:"role" validate:"required"`
}

type UpdateUserInput struct {
	NamaDepan    string `json:"namaDepan" bson:"namaDepan" form:"namaDepan" validate:"required,min=2,max=40"`
	NamaBelakang string `json:"namaBelakang" bson:"namaBelakang" form:"namaBelakang" validate:"max=60"`
	Username     string `json:"username" bson:"username" form:"username" validate:"min=5,max=30"`
	Email        string `json:"email" bson:"email" form:"email" validate:"required,email"`
	Password     string `json:"password" bson:"password" form:"password" validate:"min=6"`
	Role         string `json:"role" bson:"role" form:"role" validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required, min=1, email"`
	Password string `json:"password" validate:"required"`
}
