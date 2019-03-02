package main

import (
	"database/sql"
	"fmt"
)

func Connect() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root:12345@/test")

	if err != nil {
		return nil, fmt.Errorf("Error connect DB:__ %s", err)
	}

	return db, nil
}
