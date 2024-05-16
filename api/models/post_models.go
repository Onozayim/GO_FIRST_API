package models

import "database/sql"

type Post struct {
	Id      int64  `json:"id"`
	Post    string `json:"post"`
	User_id int64  `json:"user_id"`
}

type PostQuery struct {
	Id         *int64  `json:"id"`
	Post       *string `json:"post"`
	Created_at *string `json:"created_at"`
	Updated_at *string `json:"updated_at"`
}

type PostWithUser struct {
	Id         int64          `json:"id"`
	Post       string         `json:"post"`
	Created_at string         `json:"created_at"`
	Updated_at sql.NullString `json:"updated_at"`
	User       UserQuery      `json:"user"`
}
