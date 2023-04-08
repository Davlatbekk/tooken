package handler

import (
	"app/api/models"
	"app/config"
	"app/pkg/helper"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @ID login
// @Router /login [POST]
// @Summary Create Login
// @Description Create Login
// @Tags Login
// @Accept json
// @Produce json
// @Param Login body models.Login true "LoginRequestBody"
// @Success 201 {object} models.LoginResponse "GetLoginBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *Handler) Login(c *gin.Context) {
	var login models.Login

	err := c.ShouldBindJSON(&login)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.storages.User().GetByID(
		context.Background(),
		&models.UserPKey{Login: login.Login},
	)

	if err != nil {
		fmt.Println("okkkkkkkkkkkk")
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	if login.Password != resp.Password {
		c.JSON(http.StatusInternalServerError, errors.New("error password is not correct").Error())
		return
	}

	data := map[string]interface{}{
		"id": resp.UserId,
	}

	token, err := helper.GenerateJWT(data, config.TimeExpiredAt, h.cfg.SecretKey)
	if err != nil {
		log.Printf("error whiling GenerateJWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GenerateJWT").Error())
		return
	}

	c.JSON(http.StatusCreated, models.LoginResponse{AccessToken: token})
}

// Register godoc
// @ID register
// @Router /register [POST]
// @Summary Create Register
// @Description Create Register
// @Tags Register
// @Accept json
// @Produce json
// @Param Regester body models.Register true "RegisterRequestBody"
// @Success 201 {object} models.RegisterResponse "GetRegisterBody"
// @Response 400 {object} string "Invalid Argument"
// @Failure 500 {object} string "Server Error"
func (h *Handler) Register(c *gin.Context) {
	var register models.Register

	err := c.ShouldBindJSON(&register)
	if err != nil {
		log.Printf("error whiling create: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storages.User().Create(
		context.Background(),
		&models.Register{
			FirstName:   register.FirstName,
			LastName:    register.LastName,
			Login:       register.Login,
			Password:    register.Password,
			PhoneNumber: register.PhoneNumber,
			UserType:    register.UserType,
		},
	)

	if err != nil {
		log.Printf("error whiling GetByPKey: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GetByPKey").Error())
		return
	}

	data := map[string]interface{}{
		"id": id,
	}

	token, err := helper.GenerateJWT(data, config.TimeExpiredAt, h.cfg.SecretKey)
	if err != nil {
		log.Printf("error whiling GenerateJWT: %v\n", err)
		c.JSON(http.StatusInternalServerError, errors.New("error whiling GenerateJWT").Error())
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{AccessToken: token})
}
