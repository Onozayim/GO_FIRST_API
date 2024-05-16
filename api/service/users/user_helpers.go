package users

import (
	"api/models"
	"database/sql"
)

func ScanUserWithPost(rows *sql.Rows, user *models.UserWithPosts, post *models.PostQuery) error {
	if err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Created_at,
		&user.Updated_at,
		&post.Id,
		&post.Post,
		&post.Created_at,
		&post.Updated_at,
	); err != nil {
		return err
	}

	return nil
}
