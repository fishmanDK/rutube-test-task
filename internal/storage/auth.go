package storage

import (
	"errors"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

type AuthStorage struct {
	db *sqlx.DB
}

func NewAuthStorage(db *sqlx.DB) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

func (a *AuthStorage) IsRegisteredUser(user models.User) error {
	const op = "storage.IsRegisteredUser"

	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND hash_password = $2)"
	_, err := a.db.Exec(query, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthStorage) CreateUser(newUser models.NewUser) error {
	const op = "storage.CreateUser"

	query := "INSERT INTO users (user_name, email, hash_password, birth_date) VALUES ($1, $2, $3, $4)"
	userName := fmt.Sprintf("%s %s", newUser.LastName, newUser.FirstName)

	var birthDate string
	parts := strings.Split(newUser.BirthDate, "-")
	if len(parts) == 3 {
		year := parts[0]
		month := parts[1]
		day := parts[2]

		birthDate = fmt.Sprintf("%s-%s-%s", year, month, day)
	}

	if newUser.Surname != "" {
		userName = fmt.Sprintf("%s %s", userName, newUser.Surname)
	}
	_, err := a.db.Query(query, userName, newUser.Email, newUser.Password, birthDate)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *AuthStorage) UpdateRefreshToken(refreshToken, userEmail string) error {
	const op = "storage.UpdateRefreshToken"

	query := "UPDATE refresh_tokens SET refresh_token = $1 WHERE userID = $2"
	userId, err := getUserId(a.db, userEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.db.Query(query, refreshToken, userId)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *AuthStorage) SaveRefreshToken(refreshToken, userEmail string) error {
	const op = "storage.UpdateRefreshToken"

	userId, err := getUserId(a.db, userEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "INSERT INTO refresh_tokens (userID, refresh_token) VALUES ($1, $2)"

	_, err = a.db.Exec(query, userId, refreshToken)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *AuthStorage) CheckRefreshToken(refreshToken, userEmail string) error {
	const op = "storage.CheckRefreshToken"

	query := "SELECT EXISTS (SELECT 1 FROM users u JOIN refresh_tokens rt ON u.id = rt.userID WHERE u.id = $1 AND rt.refresh_token = $2)"

	userId, err := getUserId(a.db, userEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = a.db.Query(query, userId, refreshToken)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func getUserId(client *sqlx.DB, email string) (int, error) {
	const op = "storage.GetUserId"

	query := "SELECT id FROM users WHERE email = $1"
	var userID int
	err := client.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return -1, errors.New(fmt.Sprint(op, err.Error()))
	}
	return userID, nil
}
