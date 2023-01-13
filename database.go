package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type FileEntry struct {
	name string
	size int
	hash string
}

func createTable(db *sql.DB) error {
	sqlStmt := `CREATE TABLE IF NOT EXISTS tbl_mvdup (name TEXT, size INTEGER, hash TEXT);`
	_, err := db.Exec(sqlStmt)
	return err
}

func querySameSizeExist(db *sql.DB, size int) (bool, error) {
	rows, err := db.Query("SELECT name FROM tbl_mvdup WHERE size = $1", size)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

type asdf struct {
}

func queryByHash(db *sql.DB, hash string) (files []string, err error) {
	rows, err := db.Query("SELECT name FROM tbl_mvdup WHERE hash = $1", hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		files = append(files, name)
	}
	return files, nil
}

func insert(db *sql.DB, file FileEntry) (err error) {
	_, err = db.Exec(
		"INSERT INTO tbl_mvdup (name, size, hash) VALUES ($1, %2)",
		file.name,
		file.size,
		file.hash,
	)
	return err
}
