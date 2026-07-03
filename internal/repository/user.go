package repository

import (
	"database/sql"
	"errors"
	"messenger_v2/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserById(UserID int) (*domain.User, error) {
	var User domain.User
	err := r.db.QueryRow("SELECT name, about, avatar_url, sex FROM users WHERE id = ?", UserID).Scan(&User.Name, &User.About, &User.Avatar, &User.Sex)

	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec("UPDATE users SET name = ?, about = ? WHERE id = ?", user.Name, user.About, user.ID)
	return err
}

func (r *UserRepository) GetUserByLogin(Login string) (int, string, error) {
	var UserID int
	var PasswordFromdb string

	err := r.db.QueryRow("SELECT id, password FROM users WHERE login = ?", Login).Scan(&UserID, &PasswordFromdb)

	if err != nil {
		return 0, "", nil
	}
	return UserID, PasswordFromdb, nil
}

func (r *UserRepository) CreateUser(res *domain.RegistrationRequest, hashedPassword string) error {
	_, err := r.db.Exec("INSERT INTO users (login, password, name, sex, about, avatar_url, avatar_img) VALUES (?, ?, ?, ?, ?, ?, ?)", res.Login, hashedPassword, res.Name, res.Sex, "Пользователь TheNomax ", "unknown", "unknown")
	if err != nil {
		return errors.New("registration error")
	}
	return nil
}
