package posts

import (
	"api/models"
	"api/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

func CreateBulkPost(posts *[]models.Post, user *models.User, db *sql.DB) error {
	chunkList := utils.ChunkSlicePosts(*posts, 2)
	errorChannel := make(chan error, len(chunkList))

	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("un error inesperado ha ocurrido")
	}

	for _, chunk := range chunkList {
		go func(chunk []models.Post, tx *sql.Tx, errorChanel chan<- error) {
			valueStrings := []string{}
			valueArgs := []interface{}{}

			for _, post := range chunk {
				valueStrings = append(valueStrings, "(?, ?)")
				valueArgs = append(valueArgs, post.Post)
				valueArgs = append(valueArgs, user.Id)
			}
			stmt := fmt.Sprintf("INSERT INTO posts (post, user_id) VALUES %s", strings.Join(valueStrings, ","))

			if _, err := tx.Exec(stmt, valueArgs...); err != nil {
				log.Println("error")
				errorChannel <- errors.New("error")
				return
			}

			errorChannel <- nil
			// log.Println("error chanel is nil")
			// close(errorChanel)
		}(chunk, tx, errorChannel)
	}

	log.Println("done")

	i := 0
	for err := range errorChannel {
		if err != nil {
			log.Println("error")
			panic(err)
		} else {
			log.Println("hihihi")
			i++

			if i == len(chunkList) {
				log.Println("close")
				close(errorChannel)
			}
		}
	}

	close(errorChannel)

	tx.Commit()

	return nil
}
