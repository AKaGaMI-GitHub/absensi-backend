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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getUserCollection() *mongo.Collection {
	return config.DB.Collection("users")
}

func getUserUUID(c *gin.Context) (bson.M, string, error) {
	uuid := c.Param("uuid")
	if uuid == "" {
		return nil, "", fmt.Errorf("ID tidak boleh kosong")
	}
	return bson.M{"_id": uuid}, uuid, nil
}

func GetUsers(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursur, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		utils.ResponseJSON(c, http.StatusExpectationFailed, false, "Gagal memanggil data users!", start, err.Error())
		return
	}
	defer cursur.Close(ctx)

	var users []model.User
	if err := cursur.All(ctx, &users); err != nil {
		utils.ResponseJSON(c, http.StatusExpectationFailed, false, "Gagal memanggil data users!", start, err.Error())
		return
	}

	var userslist []model.User = []model.User{} //membuat agar default menjadi []
	for _, u := range users {
		userslist = append(userslist, model.User{
			ID:           u.ID,
			NamaDepan:    u.NamaDepan,
			NamaBelakang: u.NamaBelakang,
			Email:        u.Email,
		})
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil memanggil data users!", start, userslist)

}

func GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID Format"})
		return
	}

	var user model.User
	err = getUserCollection().FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func StoreUser(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	//check form input (Harus JSON)
	var input model.CreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Form input tidak valid!", start, err.Error())
		return
	}

	//check validasi
	if err := validate.Struct(input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Validasi gagal!", start, err.Error())
		return
	}

	//check hasing password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "Gagal hasing password!", start, err.Error())
		return
	}

	avatarPath, err := utils.UploadFile(c, "avatar", "image", "avatar")
	//check upload file
	if err != nil || avatarPath == "" {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Upload avatar gagal!", start, err.Error())
		return
	}

	role, err := GetRoleByKey(ctx, input.Role)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "Role tidak valid!", start, err.Error())
		return
	}

	newUser := model.User{
		ID:           uuid.NewString(),
		NamaDepan:    input.NamaDepan,
		NamaBelakang: input.NamaBelakang,
		Username:     input.Username,
		Email:        input.Email,
		Password:     hashedPassword,
		Avatar:       avatarPath,
		Role:         *role,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err = getUserCollection().InsertOne(ctx, newUser)
	if err != nil {
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal membuat user baru!", start, err.Error())
		return
	}

	utils.ResponseJSON(c, http.StatusCreated, true, "Berhasil membuat user baru!", start, newUser)

}

func UpdateUser(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var input model.UpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Form input tidak valid!", start, err.Error())
		return
	}

	if err := validate.Struct(input); err != nil {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Validasi gagal!", start, err.Error())
		return
	}

	id, _, err := getUserUUID(c)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "Role tidak valid!", start, err.Error())
		return
	}

	roleKey := utils.GenerateSlug(input.Role)
	if roleKey == "" {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "RoleKey tidak boleh kosong!", start, nil)
		return
	}

	updateData := bson.M{
		"$set": bson.M{
			"namaDepan":    input.NamaDepan,
			"namaBelakang": input.NamaBelakang,
			"role":         roleKey,
			"slug":         utils.GenerateSlug(input.NamaDepan, input.NamaBelakang),
		},
	}

	avatarPath, err := utils.UploadFile(c, "avatar", "image", "avatar")
	//check upload file
	if err == nil && avatarPath != "" {
		updateData["$set"].(bson.M)["avatar"] = avatarPath
	} else {
		utils.ResponseJSON(c, http.StatusUnprocessableEntity, false, "Upload avatar gagal!", start, err.Error())
		return
	}

	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			utils.ResponseJSON(c, http.StatusBadRequest, false, "Gagal hasing password!", start, err.Error())
			return
		}
		updateData["$set"].(bson.M)["password"] = hashedPassword
	}

	if input.Username != "" {
		updateData["$set"].(bson.M)["username"] = input.Username
	}

	if input.Email != "" {
		updateData["$set"].(bson.M)["email"] = input.Email
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updateUser model.User

	err = getUserCollection().FindOneAndUpdate(ctx, id, updateData, options).Decode(&updateUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.ResponseJSON(c, http.StatusNotFound, false, "Role tidak ditemukan!", start, nil)
			return
		}
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal update role!", start, err.Error())
		return
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil update user!", start, updateUser)
}

func DeleteUser(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _, err := getUserUUID(c)
	if err != nil {
		utils.ResponseJSON(c, http.StatusBadRequest, false, "UUID tidak ditemukan!", start, err.Error())
		return
	}

	result, err := getUserCollection().DeleteOne(ctx, id)
	if err != nil {
		utils.ResponseJSON(c, http.StatusInternalServerError, false, "Gagal menghapus users!", start, err.Error())
		return
	}
	if result.DeletedCount == 0 {
		utils.ResponseJSON(c, http.StatusNotFound, false, "User tidak ditemukan!", start, nil)
		return
	}

	utils.ResponseJSON(c, http.StatusOK, true, "Berhasil menghapus user!", start, nil)

}
