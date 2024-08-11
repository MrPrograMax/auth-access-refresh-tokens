package handler

import (
	"github.com/gin-gonic/gin"
	"test_task_BackDev/internal/service"
	"test_task_BackDev/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	//Для тестов
	hAuth := router.Group("/auth")
	{
		hAuth.POST("/sign-up", h.userSignUp)
		hAuth.POST("/sign-in", h.userSignIn)
	}

	//Основное задание по тз
	tokens := router.Group("/tokens")
	{
		tokens.POST("/:id", h.issueTokensPair) // Выдать Access, Refresh token по ID (UUID)
		tokens.POST("/refresh", h.userRefresh) // Обновить Access token по Refresh token
	}

	//Функция для проверки работы токенов
	test := router.Group("/test", h.userIdentity) // h.userIdentity, h.testFoo()
	{
		test.GET("", h.testFunc)
	}

	return router
}
