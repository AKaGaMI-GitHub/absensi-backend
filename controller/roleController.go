package controller

import (
	"absen-backend/config"
	"absen-backend/model"
	"absen-backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getRoleCollection() *mongo.Collection {
	return config.DB.Collection("roleUsers")
}

func GetRoleUser(c *gin.Context) {
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
