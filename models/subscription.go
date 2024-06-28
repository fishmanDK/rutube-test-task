package models

type Subscription struct {
	Email    string `db:"email"`
	UserName string `db:"user_name"`
	Birth    string `db:"birth_date"`
}
