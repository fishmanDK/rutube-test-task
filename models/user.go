package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	FirstName string `json:"fist_name"`
	LastName  string `json:"last_name"`
	Surname   string `json:"surname,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BirthDate string `json:"birth_date"`
}
