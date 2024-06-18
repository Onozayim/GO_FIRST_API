package api

import (
	"api/service/posts"
	"api/service/users"
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "api World")
	})

	mux.HandleFunc("GET /public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "", r.URL.String()))
	})

	userHandler := users.NewHandler(s.db)
	userHandler.RegisterRoutes(mux)

	postHandler := posts.NewHandler(s.db)
	postHandler.RegisterRoutes(mux)

	return http.ListenAndServe(s.addr, mux)
}
