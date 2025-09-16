package model

const (
	RoleKeyAdmin     = "admin"
	RoleKeyDosen     = "dosen"
	RoleKeyMahasiswa = "mahasiswa"
	RoleKeyPegawai   = "pegawai"
)

type Roles struct {
	ID       string `json:"id" bson:"_id"`
	RoleKey  string `bson:"roleKey" json:"roleKey"`
	RoleName string `bson:"roleName" json:"roleName"`
}

type RoleUser struct {
	RoleKey  string `bson:"roleKey" json:"roleKey"`
	RoleName string `bson:"roleName" json:"roleName"`
}

type RoleInput struct {
	RoleName string `json:"roleName" form:"roleName" validate:"required,min=5"`
}