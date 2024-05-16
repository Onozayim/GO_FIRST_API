package models

import "database/sql"

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserQuery struct {
	Id         int64          `json:"id"`
	Username   string         `json:"username"`
	Email      string         `json:"email"`
	Created_at string         `json:"created_at"`
	Updated_at sql.NullString `json:"updated_at"`
}

type UserWithPosts struct {
	Id         int64        `json:"id"`
	Username   string       `json:"username"`
	Email      string       `json:"email"`
	Created_at string       `json:"created_at"`
	Updated_at *string      `json:"updated_at"`
	Posts      []*PostQuery `json:"posts"`
}
