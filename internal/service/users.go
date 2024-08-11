package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"test_task_BackDev/internal/domain"
	"test_task_BackDev/internal/repository"
	"test_task_BackDev/pkg/auth"
	"test_task_BackDev/pkg/email"
	"time"
)

type UsersService struct {
	repo         repository.Users
	tokenManager auth.TokenManager

	emailService email.Sender

	accessTokenTLL  time.Duration
	refreshTokenTLL time.Duration
}

func NewUsersService(repo repository.Users, tokenManager auth.TokenManager, emailService email.Sender, accessTLL, refreshTLL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		tokenManager:    tokenManager,
		emailService:    emailService,
		accessTokenTLL:  accessTLL,
		refreshTokenTLL: refreshTLL,
	}
}

func (s *UsersService) SignUp(input UserInput) (uuid.UUID, error) {
	passwordHash, err := generatePasswordHash(input.Password)
	if err != nil {
		return uuid.Nil, err
	}

	id := uuid.New()

	user := domain.User{
		ID:       id,
		Email:    input.Email,
		Password: passwordHash,
		Ip:       input.Ip,
	}

	return s.repo.Create(user)
}

func (s *UsersService) SignIn(input UserInput) (Tokens, error) {
	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return Tokens{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return Tokens{}, errors.New("wrong password")
	}

	return s.createSession(user.ID, input.Ip)
}

func (s *UsersService) RefreshToken(refreshToken, userIp string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	if user.ExpiresAt.Before(time.Now()) {
		return Tokens{}, errors.New("token expired")
	}

	if user.Ip != userIp {
		logrus.Warning(fmt.Sprintf("Suspicious activity. User[%s] used NEW IP(%s). Old IP(%s)]", user.ID, userIp, user.Ip))
		warningEmail := email.GenerateIpWarningEmail(user.Email, userIp)
		err = s.emailService.Send(warningEmail)
		if err != nil {
			logrus.Warning("Mock email send")
		} else {
			logrus.Info("mail sent to ", user.Email)
		}
	}

	return s.createSession(user.ID, userIp)
}

func (s *UsersService) IssueTokensPair(userId uuid.UUID) (Tokens, error) {
	user, err := s.repo.GetById(userId)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(user.ID, user.Ip)
}

func (s *UsersService) createSession(userId uuid.UUID, ip string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.String(), ip, s.accessTokenTLL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = generatePasswordHash(res.RefreshToken)

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		IpAddress:    ip,
		ExpiresAt:    time.Now().Add(s.refreshTokenTLL),
	}

	err = s.repo.SetSession(userId, session)

	return res, err
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("invalid bcrypt hash generation")
	}

	return string(hash), nil
}
