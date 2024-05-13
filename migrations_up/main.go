package main

import (
	"log"
	migration_model "migrations_up/migrations"
)

func main() {

	migration_handler, err := migration_model.GET_MIGRATIONS()

	if err != nil {
		log.Fatal(err)
	}

	migration_handler.Migrate("create_user", `
		CREATE TABLE IF NOT EXISTS users(
	 		id int(11) NOT NULL auto_increment,
	 		username varchar(50) not null,
	 		email varchar(100) not null,
	 		password varchar(100) not null,
	 		create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	 		update_at DATETIME ON UPDATE CURRENT_TIMESTAMP,
			
	 		PRIMARY KEY (id)
	 	)
	`)

	migration_handler.Migrate("create_post", `
		CREATE TABLE IF NOT EXISTS posts(
			id int(11) NOT NULL auto_increment,
			user_id int(11) NOT NULL,
			post varchar(255) NOT NULL,
			create_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			update_at DATETIME ON UPDATE CURRENT_TIMESTAMP,
		
			PRIMARY KEY (id),
			CONSTRAINT FK_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
		)
	`)
}
