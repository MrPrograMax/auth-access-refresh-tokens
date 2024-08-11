package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"test_task_BackDev/internal/service"
)

type userInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Ip       string `json:"ip"`
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) userSignUp(c *gin.Context) {
	var inp userInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Users.SignUp(service.UserInput{
		Email:    inp.Email,
		Password: inp.Password,
		Ip:       c.ClientIP(),
	})

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user already exists")
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"Id": id})
}

func (h *Handler) userSignIn(c *gin.Context) {
	var inp userInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	res, err := h.services.Users.SignIn(service.UserInput{
		Email:    inp.Email,
		Password: inp.Password,
		Ip:       c.ClientIP(),
	})

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

func (h *Handler) issueTokensPair(c *gin.Context) {
	userId := uuid.MustParse(c.Param("id"))

	res, err := h.services.Users.IssueTokensPair(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

func (h *Handler) userRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	userIp := c.ClientIP()
	res, err := h.services.Users.RefreshToken(inp.Token, userIp)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AccessToken":  res.AccessToken,
		"RefreshToken": res.RefreshToken,
	})
}

func (h *Handler) testFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "You are authorized!",
	})
}
