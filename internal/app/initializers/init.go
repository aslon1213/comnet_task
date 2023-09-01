package initializers

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Init function prepares db connection and does migrations
func Init() (*sql.DB, error) {
	// fmt.Println(os.Getwd())
	// if dir, _ := os.Getwd(); dir != "/home/aslon/go/src/github.com/aslon1213/comnet_task" {

	// }
	// set the timezone to UTC+5
	// os.Setenv("TZ", "Asia/Tashkent")

	// load envs
	err := load_envs()
	if err != nil {
		return nil, err
	}

	if gin.Mode() == gin.TestMode {
		err := os.Chdir("/Users/aslonkhamidov/Desktop/code/tasks/comnet_task/")
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite3", "./db/db.sqlite3")

	if err != nil {
		return nil, err
	}
	createtables(db)
	return db, nil
}

// do migrations
func createtables(db *sql.DB) {

	file, err := os.ReadFile("db/migrations/migrations.sql")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(string(file))
	if err != nil {
		panic(err)
	}

}
