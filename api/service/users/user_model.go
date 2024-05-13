package users

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPayload struct {
	Id         int64          `json:"id"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Created_at string         `json:"created_at"`
	Updated_at sql.NullString `json:"updated_at"`
}

func CreateUser(user *User, db *sql.DB) error {

	if user.Username == "" {
		return fmt.Errorf("username is empty")
	}

	if user.Email == "" {
		return fmt.Errorf("email is empty")
	}

	if user.Email == "" {
		return fmt.Errorf("password is empty")
	}

	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, user.Username, user.Email, user.Password)

	if err != nil {
		return err
	}

	id, err := insertResult.LastInsertId()

	if err != nil {
		return err
	}

	user.Id = id

	return nil
}
