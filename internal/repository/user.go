package repository

import (
	"database/sql"
	"errors"
	"log"
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
	err := r.db.QueryRow("SELECT id, name, about, avatar_url, sex, COALESCE(group_name, '') FROM users WHERE id = ?", UserID).Scan(&User.ID, &User.Name, &User.About, &User.Avatar, &User.Sex, &User.Group)

	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec("UPDATE users SET name = ?, about = ?, group_name = ? WHERE id = ?", user.Name, user.About, user.Group, user.ID)
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
	_, err := r.db.Exec("INSERT INTO users (login, password, name, sex, about, avatar_url, avatar_img, group_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", res.Login, hashedPassword, res.Name, res.Sex, "Пользователь TheNomax ", "unknown", "unknown", res.Group)
	if err != nil {
		return errors.New("registration error")
	}
	return nil
}

func (r *UserRepository) UploadAvatarUser(UserID int, AvatarURL string) error {
	_, err := r.db.Exec("UPDATE users SET avatar_url = ? WHERE id = ? ", AvatarURL, UserID)
	if err != nil {
		log.Println("БД: Ошибка при обнавлении данных пользователя")
		return err
	}
	return nil
}

func (r *UserRepository) GetPersonByID(PersonID int) (*domain.UserResponse, error) {
	var Person domain.UserResponse
	err := r.db.QueryRow("SELECT name, about, avatar_url, sex, COALESCE(group_name, '') FROM users WHERE id = ?", PersonID).Scan(&Person.Name, &Person.About, &Person.Avatar, &Person.Sex, &Person.Group)
	if err != nil {
		log.Println("БД: Ошибка при получении данных другого пользователя")
		return nil, err
	}
	return &Person, nil
}
