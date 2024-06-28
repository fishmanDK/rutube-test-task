package service

import (
	"github.com/fishmanDK/rutube-test-task/internal/storage"
	"github.com/fishmanDK/rutube-test-task/models"
)

type Service struct {
	Auth
	Api
}

type Auth interface {
	Authentication(user models.User) (models.Tokens, error)
	CreateUser(newUser models.NewUser) error
	ParseToken(accessToken string) (string, error)
	ApdateAccessToken(updateAccessTokenReq models.UpdateRefreshTokenRequest) (models.Tokens, error)
}

type Api interface {
	Subscribe(rootEmail, subsEmail string) error
	Unsubscribe(rootEmail, subsEmail string) error
	GetAllSubs(userEmail string) map[string]models.Subscription
}

func MustService(storage *storage.Storage) *Service {
	return &Service{
		Auth: NewAuthService(storage),
		Api:  NewApiService(storage),
	}
}
