package service

import (
	"github.com/google/uuid"
	"test_task_BackDev/internal/repository"
	"test_task_BackDev/pkg/auth"
	"test_task_BackDev/pkg/email"
	"time"
)

type UserInput struct {
	Email    string
	Password string
	Ip       string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(input UserInput) (uuid.UUID, error)
	SignIn(input UserInput) (Tokens, error)
	RefreshToken(refreshToken, userIp string) (Tokens, error)
	IssueTokensPair(uuid uuid.UUID) (Tokens, error)
}

type Services struct {
	Users
}

type Deps struct {
	Repos           *repository.Repository
	TokenManager    auth.TokenManager
	EmailSender     email.Sender
	AccessTokenTLL  time.Duration
	RefreshTokenTLL time.Duration
}

func NewService(deps Deps) *Services {
	return &Services{
		Users: NewUsersService(deps.Repos, deps.TokenManager, deps.EmailSender, deps.AccessTokenTLL, deps.RefreshTokenTLL),
	}
}
