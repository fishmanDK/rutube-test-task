package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/config"
	"github.com/fishmanDK/rutube-test-task/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Auth interface {
	IsRegisteredUser(user models.User) error
	CreateUser(newUser models.NewUser) error
	UpdateRefreshToken(refreshToken, userEmail string) error
	SaveRefreshToken(refreshToken, userEmail string) error
	CheckRefreshToken(refreshToken, userEmail string) error
}

type Api interface {
	Subscribe(rootEmail, subsEmail string) error
	Unsubscribe(rootEmail, subsEmail string) error
	GetAllSubs(userEmail string) map[string]models.Subscription
}

type Storage struct {
	Auth
	Api
}

func MustStorage(cfg *config.ConfigDB) (*Storage, error) {
	const op = "storage.MustStorage"
	fmt.Println(cfg.ToString())
	client, err := sql.Open("postgres", cfg.ToString())

	if err != nil {
		return nil, errors.New(fmt.Sprint(op, err))
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	}

	sqlxDB := sqlx.NewDb(client, "postgres")

	return &Storage{
		Auth: NewAuthStorage(sqlxDB),
		Api:  NewApiStorage(sqlxDB),
	}, nil
}
