package main

import (
	"api/cmd/api"
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
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
	server := api.NewApIServer(":8080", db)

	log.Printf("Server runing on: http://localhost:8080/ \n")

	if err := server.Run(); err != nil {
		log.Fatal("ERROR")
	}
}
