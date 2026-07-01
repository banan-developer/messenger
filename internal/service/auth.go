package service

import (
	"errors"
	"fmt"
	"messenger_v2/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (a *AuthService) Login(login string, password string) (int, error) {
	UserID, PasswordFromdb, err := a.userRepo.GetUserByLogin(login)

	if err != nil {
		return 0, errors.New("Invalid credentials")
	}

	HashError := bcrypt.CompareHashAndPassword(
		[]byte(PasswordFromdb),
		[]byte(password),
	)
	if HashError != nil {
		return 0, errors.New("Ошибка авторизации")
	}
	return UserID, nil
}

func (a *AuthService) Registration(login, password, name, sex string) error {
	hashPassword, _ := hashPassword(password)
	about := "Пользователь TheNomax"
	avatar_url := "unknown"
	avatar_img := "unknow"

	err := a.userRepo.CreateUser(login, hashPassword, name, sex, about, avatar_url, avatar_img)
	if err != nil {
		return fmt.Errorf("ошибка при создании пользователя %w", err)
	}
	return nil
}
