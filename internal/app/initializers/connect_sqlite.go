package initializers

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlite() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "db/db.sqlite3")
	if err != nil {
		return nil, err
	}
	// fmt.Println(db.Conn(context.Background()))
	fmt.Println("Connected to sqlite3 database")
	// defer db.Close()
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
