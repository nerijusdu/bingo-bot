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
	db, err := sql.Open("sqlite3", "./data/data.db")
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

func (d *Database) Exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err.Error())
	}

	return err
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	res, err := d.db.Query(query, args...)
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err.Error())
	}

	return res, err
}

func (d *Database) Insert(query string, args ...interface{}) (int, error) {
	res, err := d.db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error executing query: %s\n", err.Error())
	}

	return int(id), err
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
