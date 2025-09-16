package model

import (
	"time"
)

const (
	RoleKeyAdmin     = "admin"
	RoleKeyDosen     = "dosen"
	RoleKeyMahasiswa = "mahasiswa"
	RoleKeyPegawai   = "pegawai"
)

type UserRole string

type User struct {
	ID           string    `json:"id" bson:"_id"`
	NamaDepan    string    `json:"namaDepan" bson:"namaDepan"`
	NamaBelakang string    `json:"namaBelakang" bson:"namaBelakang"`
	Email        string    `json:"email" bson:"email"`
	Password     string    `json:"-" bson:"password"` // jangan expose ke JSON
	Avatar       string    `json:"avatar" bson:"avatar"`
	Role         RoleUser  `json:"role" bson:"role"`
	Slug         string    `json:"slug" bson:"slug"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updated_at"`
}

type RoleUser struct {
	RoleKey  string `bson:"roleKey" json:"roleKey"`
	RoleName string `bson:"roleName" json:"roleName"`
}

type CreateUserInput struct {
	NamaDepan    string `json:"namaDepan" form:"namaDepan" validate:"required,min=2,max=40"`
	NamaBelakang string `json:"namaBelakang" form:"namaBelakang" validate:"max=60"`
	Email        string `json:"email" form:"email" validate:"required,email"`
	Password     string `json:"password" form:"password" validate:"required,min=6"`
}

type UpdateUserInput struct {
	NamaDepan    string `json:"namaDepan" bson:"namaDepan" validate:"required, min=2, max=40"`
	NamaBelakang string `json:"namaBelakang" bson:"namaBelakang" validate:"max=60"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required, min=1, email"`
	Password string `json:"password" validate:"required"`
}
