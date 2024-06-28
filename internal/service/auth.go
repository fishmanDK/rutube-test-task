package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/internal/storage"
	"github.com/fishmanDK/rutube-test-task/models"
)

type AuthService struct {
	repo *storage.Storage
}

func NewAuthService(repo *storage.Storage) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (a *AuthService) Authentication(user models.User) (models.Tokens, error) {
	const op = "service.Authentication"

	user.Password = HashPassword(user.Password)

	err := a.repo.Auth.IsRegisteredUser(user)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := CreateAccessToken(user.Email)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := CreateRefreshToken()
	if err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := a.repo.Auth.SaveRefreshToken(refreshToken, user.Email); err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func (a *AuthService) CreateUser(newUser models.NewUser) error {
	const op = "service.CreateUser"

	newUser.Password = HashPassword(newUser.Password)
	err := a.repo.Auth.CreateUser(newUser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) ApdateAccessToken(updateAccessTokenReq models.UpdateRefreshTokenRequest) (models.Tokens, error) {
	const op = "service.ApdateAccessToken"

	if err := a.repo.Auth.CheckRefreshToken(updateAccessTokenReq.RefreshToken, updateAccessTokenReq.Email); err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := CreateAccessToken(updateAccessTokenReq.Email)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	newRefreshToken, err := CreateRefreshToken()
	if err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := a.repo.Auth.UpdateRefreshToken(newRefreshToken, updateAccessTokenReq.Email); err != nil {
		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	return tokens, nil
}

func HashPassword(password string) string {
	data := []byte(password + salt)
	hashData := sha256.Sum256(data)
	hashString := hex.EncodeToString(hashData[:])

	return hashString
}
