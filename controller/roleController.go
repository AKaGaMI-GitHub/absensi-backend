package controller

import (
	"absen-backend/config"
	"absen-backend/model"
	"absen-backend/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//private
func getRoleCollection() *mongo.Collection {
	return config.DB.Collection("roleusers")
}

//private
func getRoleUUID(c *gin.Context) (bson.M, string, error) {
	uuid := c.Param("uuid")
	if uuid == "" {
		return nil, "", fmt.Errorf("ID tidak boleh kosong")
	}
	return bson.M{"_id": uuid}, uuid, nil
}

//public
func GetRoleUser(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursur, err := getRoleCollection().Find(ctx, bson.M{})
	if err != nil {
		utils.ResponseJSON(c, http.StatusExpectationFailed, false, "Gagal memanggil data roles!", start, err.Error())
		return
	}
	defer cursur.Close(ctx)

	var roles []model.Roles
	if err := cursur.All(ctx, &roles); err != nil {
		utils.ResponseJSON(c, http.StatusExpectationFailed, false, "Gagal memanggil data roles!", start, err.Error())
		return
	}

	var roleList []model.Roles = []model.Roles{} //membuat agar default menjadi []
	for _, u := range roles {
		roleList = append(roleList, model.Roles{
			RoleKey:    u.RoleKey,
			RoleName: u.RoleName,
		})
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil memanggil data roles!", start, roleList)
}

func GetRoleByKey(ctx context.Context, key string) (*model.RoleUser, error) {
	if key == "" {
		return nil, fmt.Errorf("role key kosong")
	}

	var role model.RoleUser
	filter := bson.M{"roleKey": key}
	err := getRoleCollection().FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("role dengan key %s tidak ditemukan", key)
		}
		return nil, err
	}

	return &role, nil
}

//public
func StoreRole(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var input model.RoleInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Form input tidak valid!", start, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Validasi gagal!", start, err.Error())
		return
	}

	roleKey := utils.GenerateSlug(input.RoleName)
	fmt.Println("RoleName:", input.RoleName, "=> RoleKey:", roleKey)
	if roleKey == "" {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "RoleKey tidak boleh kosong!", start, nil)
		return
	}
	newRoles := model.Roles{
		ID:           uuid.NewString(),
		RoleKey:   roleKey,
		RoleName: input.RoleName,
	}
	
	_, err := getRoleCollection().InsertOne(ctx, newRoles)
	if err != nil {
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal membuat role baru!", start, err.Error())
		return
	}

	utils.ResponseJSON(c, http.StatusCreated, true, "Berhasil membuat role baru!", start, newRoles)
}

//public
func UpdateRole(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var input model.RoleInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Form input tidak valid!", start, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Validasi gagal!", start, err.Error())
		return
	}

	//find data dengan id
	filterRoles, _, err := getRoleUUID(c)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, err.Error(), start, nil)
		return
	}

	roleKey := utils.GenerateSlug(input.RoleName)
	if roleKey == "" {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "RoleKey tidak boleh kosong!", start, nil)
		return
	}

	//Update data!!
	updateData := bson.M{
		"$set": bson.M{
			"roleKey":   roleKey,
			"roleName":  input.RoleName,
		},
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedRole model.Roles

    // Find and update the document, and decode the result into updatedRole
    err = getRoleCollection().FindOneAndUpdate(ctx, filterRoles, updateData, options).Decode(&updatedRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ResponseJSON(c, http.StatusNotFound, false, "Role tidak ditemukan!", start, nil)
			return
		}
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal update role!", start, err.Error())
		return
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil update role!", start, updatedRole)
}

//public
func DeleteRole(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	filterRoles, _, err := getRoleUUID(c)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "UUID tidak ditemukan!", start, err.Error())
		return
	}
	
	result, err := getRoleCollection().DeleteOne(ctx, filterRoles)
	if err != nil {
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal menghapus role!", start, err.Error())
		return
	}

	if result.DeletedCount == 0 {
		utils.ResponseJSON(c, http.StatusNotFound, false, "Role tidak ditemukan!", start, nil)
		return
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil menghapus role!", start, nil)
}
