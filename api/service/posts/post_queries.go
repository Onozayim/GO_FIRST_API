package posts

import (
	"api/models"
	"context"
	"database/sql"
	"fmt"
)

func CreatePost(post *models.Post, user *models.User, db *sql.DB) error {

	if post.Post == "" {
		return fmt.Errorf("post vacio")
	}

	query := "INSERT INTO posts (post, user_id) VALUES (?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, post.Post, user.Id)

	if err != nil {
		return err
	}

	id, err := insertResult.LastInsertId()

	if err != nil {
		return err
	}

	post.Id = id

	return nil
}

func GetPost(user_id int64, db *sql.DB) (models.PostWithUser, error) {
	query := `SELECT
				p.id as id,
				p.post,
				p.create_at as p_created_at,
				p.update_at as p_updated_at,
				u.id as u_id,
				u.username,
				u.email,
				u.create_at as u_created_at,
				u.update_at as u_updated_at
			FROM
				posts as p
				INNER JOIN users as u ON u.id = p.user_id
			WHERE p.id = ?`

	post := models.PostWithUser{}

	err := db.QueryRow(query, user_id).Scan(
		&post.Id,
		&post.Post,
		&post.Created_at,
		&post.Updated_at,
		&post.User.Id,
		&post.User.Username,
		&post.User.Email,
		&post.User.Created_at,
		&post.User.Updated_at,
	)

	if err != nil {
		return post, err
	}

	return post, err
}
