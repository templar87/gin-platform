package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"strconv"
	"hzl.im/gin-platform/services"
	"hzl.im/gin-platform/models"
)

func UserList(ctx *gin.Context) {
	users := []models.User{}

	offset, err := strconv.Atoi(ctx.Param("offset"))
	if err!=nil {
		offset = 0
	}

	limit, err := strconv.Atoi(ctx.Param("limit"))
	if err!=nil {
		limit = 20
	}

	rows := services.DB.Order("id desc").Limit(limit).Offset(offset).Find(&users)
	if rows.Error != nil {
		res := models.ResultData{9998, "list user fail", ""}
		ctx.JSON(http.StatusInternalServerError, res)
	}

	res := models.ResultData{0, "", rows.Value}
	ctx.JSON(http.StatusOK, res)
}


func UserInfo(ctx *gin.Context) {
	user_id := ctx.Param("user_id")
	user := &models.User{}
	if err := services.DB.First(&user, user_id).Error; err != nil {
		res := models.ResultData{1, "user not exist", ""}
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := models.ResultData{0, "", user}
	ctx.JSON(http.StatusOK, res)
}

func UserAdd(ctx *gin.Context) {
	user := &models.User{}
	if ctx.Bind(user) != nil {
		res := models.ResultData{9999, "params wrong", ""}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := services.DB.Create(user).Error; err != nil {
		res := models.ResultData{9998, "add user fail", ""}
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := models.ResultData{0, "", ""}
	ctx.JSON(http.StatusOK, res)
}

func UserDel(ctx *gin.Context) {
	user_id := ctx.PostForm("user_id")
	if _, err := strconv.Atoi(user_id); err!=nil {
		res := models.ResultData{9999, "params wrong", ""}
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user := &models.User{}
	if err := services.DB.First(&user, user_id).Error; err != nil {
		res := models.ResultData{1, "user not exist", ""}
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if err := services.DB.Delete(&user).Error; err != nil {
		res := models.ResultData{9998, "del user fail", ""}
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := models.ResultData{0, "", ""}
	ctx.JSON(http.StatusOK, res)
}