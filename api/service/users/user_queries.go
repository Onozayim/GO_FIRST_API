package users

import (
	"api/models"
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

func CreateUser(user *models.User, db *sql.DB) error {
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

func GetUserPosts(user_id int64, db *sql.DB) (*models.UserWithPosts, error) {
	user := models.UserWithPosts{}
	var err error

	query := `
		SELECT
			u.id,
			u.username,
			u.email,
			u.create_at,
			u.update_at,
			p.id as post_id,
			p.post,
			p.create_at as p_created_at,
			p.update_at as p_updated_at
		FROM
			users as u
			LEFT JOIN posts as p on p.user_id = u.id
		where
			u.id = ?
	`

	rows, err := db.Query(query, user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		post := &models.PostQuery{}
		if err = ScanUserWithPost(rows, &user, post); err != nil {
			return nil, nil
		}

		if post.Id != nil {
			user.Posts = append(user.Posts, post)
		}
		post = nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(user, models.UserWithPosts{}) {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

func GetAllUserPosts(db *sql.DB) ([]*models.UserWithPosts, error) {
	usermap := make(map[int64]*models.UserWithPosts)
	var err error

	query := `
		SELECT
			u.id,
			u.username,
			u.email,
			u.create_at,
			u.update_at,
			p.id as post_id,
			p.post,
			p.create_at as p_created_at,
			p.update_at as p_updated_at
		FROM
			users as u
			LEFT JOIN posts as p on p.user_id = u.id
	`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := &models.UserWithPosts{}
		post := &models.PostQuery{}

		if err = ScanUserWithPost(rows, user, post); err != nil {
			return nil, err
		}

		new_user, ok := usermap[user.Id]

		if !ok {
			new_user = user
			usermap[user.Id] = new_user
			new_user.Posts = nil
		}

		if post.Id != nil {
			new_user.Posts = append(new_user.Posts, post)
		}

		user = nil
		post = nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(usermap) == 0 {
		return nil, fmt.Errorf("empty result")
	}

	data := make([]*models.UserWithPosts, 0, len(usermap))

	for _, u := range usermap {
		data = append(data, u)
	}

	return data, nil
}
