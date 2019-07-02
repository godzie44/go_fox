package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/user/godzie44/go_fox/fixture"
	"github.com/user/godzie44/go_fox/search"
	"log"
	"net/http"
)

var (
	pgUser = flag.String("pg_user", "postgres", "Pgsql user name")

	pgPw = flag.String("pg_password", "postgres", "Pgsql user password")

	pgHost = flag.String("pg_host", "localhost:5432", "Pgsql host")

	pgDb = flag.String("pg_db", "postgres", "Pgsql database name")
)

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main()  {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", *pgUser, *pgPw, *pgHost, *pgDb)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	fixture.GenerateData(db)

	db.SetMaxOpenConns(20)

	searchEngine := search.NewSearchEngine(db)
	r := mux.NewRouter()

	r.Handle("/{id1:[0-9]+}/{id2:[0-9]+}", panicMiddleware(searchEngine)).Methods("GET")

	log.Println("starting server at :12345")
	log.Fatal(http.ListenAndServe(":12345", r))
}
