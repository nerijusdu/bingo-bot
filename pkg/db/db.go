package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "./data.db")
	checkErr(err)

	checkErr(db.Ping())

	initializeDatabase(db)

	fmt.Println("Database ready")

	return &Database{
		db: db,
	}
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) Exec(query string, args ...interface{}) sql.Result {
	res, err := d.db.Exec(query, args...)
	checkErr(err)

	return res
}

func (d *Database) Query(query string, args ...interface{}) *sql.Rows {
	rows, err := d.db.Query(query, args...)
	checkErr(err)

	return rows
}

func initializeDatabase(db *sql.DB) {
	path := filepath.Join("./", "db", "initialize.sql")
	c, err := ioutil.ReadFile(path)
	checkErr(err)

	sql := string(c)
	_, err = db.Exec(sql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
