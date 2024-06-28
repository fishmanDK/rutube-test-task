package storage

import (
	"errors"
	"fmt"
	"github.com/fishmanDK/rutube-test-task/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type ApiStorage struct {
	db *sqlx.DB
}

func NewApiStorage(db *sqlx.DB) *ApiStorage {
	return &ApiStorage{
		db: db,
	}
}

func (a *ApiStorage) Subscribe(rootEmail, subsEmail string) error {
	const op = "storage.Subscribe"

	idRootEmail, err := getUserId(a.db, rootEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	idSubsEmail, err := getUserId(a.db, subsEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "INSERT INTO birth_subs (root_user, subscriber) VALUES ($1, $2)"
	_, err = a.db.Exec(query, idRootEmail, idSubsEmail)
	if err != nil {
		return errors.New(fmt.Sprint(op, err.Error()))
	}
	return nil
}

func (a *ApiStorage) Unsubscribe(rootEmail, subsEmail string) error {
	const op = "storage.Unsubscribe"

	idRootEmail, err := getUserId(a.db, rootEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	idSubsEmail, err := getUserId(a.db, subsEmail)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := "DELETE FROM birth_subs WHERE root_user = $1 AND subscriber = $2"
	_, err = a.db.Exec(query, idRootEmail, idSubsEmail)
	if err != nil {
		return errors.New(fmt.Sprint(op, err.Error()))
	}

	return nil
}

func (a *ApiStorage) GetAllSubs(userEmail string) map[string]models.Subscription {
	const op = "storage.GetAllSubs"

	tx, err := a.db.Begin()
	defer tx.Rollback()

	subs := make(map[string]models.Subscription)
	query := "SELECT u.user_name, u.email, u.birth_date FROM users u JOIN birth_subs bs ON u.id = bs.root_user WHERE bs.subscriber = $1;"
	userId, err := getUserId(a.db, userEmail)
	if err != nil {
		return subs
	}
	rows, err := tx.Query(query, userId)
	if err != nil {
		return subs
	}
	defer rows.Close()
	for rows.Next() {
		sub := models.Subscription{}
		err := rows.Scan(&sub.UserName, &sub.Email, &sub.Birth)
		if err != nil {
			return subs
		}

		parsedDate, err := time.Parse(time.RFC3339, sub.Birth)
		if err != nil {
			fmt.Println("Ошибка при разборе даты:", err)
			return subs
		}
		formattedDate := parsedDate.Format("02.01.2006")
		sub.Birth = formattedDate
		subs[sub.Email] = sub
	}

	return subs
}
