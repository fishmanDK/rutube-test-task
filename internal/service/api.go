package service

import (
	"errors"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/internal/storage"
	"github.com/fishmanDK/rutube-test-task/models"
)

type ApiService struct {
	repo *storage.Storage
}

func NewApiService(repo *storage.Storage) *ApiService {
	return &ApiService{
		repo: repo,
	}
}

func (a *ApiService) Subscribe(rootEmail, subsEmail string) error {
	const op = "service.Subscribe"

	if err := a.repo.Api.Subscribe(rootEmail, subsEmail); err != nil {
		return errors.New(fmt.Sprintf("%s %s", op, err.Error()))
	}

	return nil
}

func (a *ApiService) Unsubscribe(rootEmail, subsEmail string) error {
	const op = "service.Unsubscribe"

	if err := a.repo.Api.Unsubscribe(rootEmail, subsEmail); err != nil {
		return errors.New(fmt.Sprintf("%s %s", op, err.Error()))
	}

	return nil
}

func (a *ApiService) GetAllSubs(userEmail string) map[string]models.Subscription {
	subs := a.repo.Api.GetAllSubs(userEmail)
	return subs
}
