package migration_model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type MIGRATION_STRUCT struct {
	migration      string
	migration_name string
}

type MIGRATION_HANDLER struct {
	names map[string]bool
	db    *sql.DB
}

func NEW_MIGRATION_STRUCT(migration_name string, migration string) *MIGRATION_STRUCT {
	return &MIGRATION_STRUCT{
		migration:      migration,
		migration_name: migration_name,
	}
}

func GET_MIGRATIONS() (*MIGRATION_HANDLER, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "go_tutorial",
		AllowNativePasswords: true,
	}

	var err error

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	if pingErr := db.Ping(); pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Conectado!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id int(11) auto_increment NOT NULL,
			migration_name varchar(50) NOT NULL, 
			UNIQUE(migration_name),
			PRIMARY KEY (id)
		)
	`)

	if err != nil {
		fmt.Println("ERROR CREATING MIGRATION ERROR")
		return nil, nil
	}

	rows, err := db.Query("SELECT migration_name FROM migrations")

	var migration_names = make(map[string]bool)

	if err != nil {
		return &MIGRATION_HANDLER{names: migration_names, db: db}, err
	}

	defer rows.Close()

	for rows.Next() {
		var migration_name string

		if err = rows.Scan(&migration_name); err != nil {
			return &MIGRATION_HANDLER{names: migration_names, db: db}, err
		}

		migration_names[migration_name] = true
	}

	if err = rows.Err(); err != nil {
		return &MIGRATION_HANDLER{names: migration_names, db: db}, err
	}

	return &MIGRATION_HANDLER{
		names: migration_names,
		db:    db,
	}, nil
}

func (s *MIGRATION_HANDLER) Migrate(migration_name string, migration string) error {
	if s.names[migration_name] {
		log.Printf("%s Alredy migrated \n", migration_name)
		return nil
	}

	if utf8.RuneCountInString(migration_name) > 50 {
		log.Printf("%s cant have more than 50 characters \n", migration_name)

		return nil
	}

	log.Printf("migrating: %s...", migration_name)

	_, err := s.db.Exec(migration)

	if err != nil {
		log.Fatal("ERROR MIGRATING: " + migration_name)
		return fmt.Errorf("error: %v", err)
	}

	_, err = s.db.Exec("INSERT INTO migrations (migration_name) values (?)", migration_name)

	if err != nil {
		log.Fatal("ERROR MIGRATING: " + migration_name)
		return fmt.Errorf("error: %v", err)
	}

	log.Printf("%s migrated", migration_name)

	return nil
}
