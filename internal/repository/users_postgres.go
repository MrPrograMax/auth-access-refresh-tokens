package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"test_task_BackDev/internal/domain"
)

type UsersRepo struct {
	db *sqlx.DB
}

func NewUsersRepo(db *sqlx.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Create(user domain.User) (uuid.UUID, error) {
	var id uuid.UUID

	insertQuery := fmt.Sprintf(`INSERT INTO %s (id, email, password, ip) VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)
	row := r.db.QueryRow(insertQuery, user.ID, user.Email, user.Password, user.Ip)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *UsersRepo) GetByRefreshToken(refreshToken string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE refresh_token = $1`, usersTable)
	err := r.db.Get(&user, query, refreshToken)

	return user, err
}

func (r *UsersRepo) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT u.id, u.email, u.password FROM %s u WHERE email = $1`, usersTable)
	err := r.db.Get(&user, query, email)

	return user, err
}

func (r *UsersRepo) SetSession(userId uuid.UUID, session domain.Session) error {
	query := fmt.Sprintf(`UPDATE %s SET refresh_token = $1, expires_at = $2, ip = $3 WHERE id = $4`, usersTable)
	_, err := r.db.Exec(query, session.RefreshToken, session.ExpiresAt, session.IpAddress, userId)
	return err
}

func (r *UsersRepo) GetById(userId uuid.UUID) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, usersTable)
	err := r.db.Get(&user, query, userId)

	return user, err
}
