package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	for i := 0; i < 10; i++ {
		db, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}
	}
	return nil, fmt.Errorf("could not connect to database")
}
