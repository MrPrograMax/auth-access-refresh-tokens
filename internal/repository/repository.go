package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"test_task_BackDev/internal/domain"
)

type Users interface {
	Create(user domain.User) (uuid.UUID, error)
	GetByRefreshToken(refreshToken string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	SetSession(userId uuid.UUID, session domain.Session) error
	GetById(userId uuid.UUID) (domain.User, error)
}

type Repository struct {
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUsersRepo(db),
	}
}
