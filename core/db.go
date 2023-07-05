package core

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	Session *sql.DB
	Context context.Context
	stop    context.CancelFunc
}

var Connection *DBConnection

func OpenDB() {
	// Open the database file in read-write mode.
	var err error
	Connection = new(DBConnection)

	Connection.Session, err = sql.Open("sqlite3", "data/fake.db")
	if err != nil {
		panic(err)
	}

	// Check if the database exists.
	if _, err := os.Stat("data/fake.db"); os.IsNotExist(err) {
		panic("Database does not exist.")
	}

	Connection.Session.SetConnMaxLifetime(0)
	Connection.Session.SetMaxIdleConns(3)
	Connection.Session.SetMaxOpenConns(3)

	Connection.Context, Connection.stop = context.WithCancel(context.Background())
	// Do something with the database.
	fmt.Println("Database opened successfully.")
}

func CloseDB() {
	Connection.Session.Close()
	Connection.stop()
}

func Get(query string, params ...string) (*sql.Row, error) {
	stmt, err := Connection.Session.Prepare(query)
	if err != nil {
		return nil, err
	} else {
		defer stmt.Close()
	}
	result := stmt.QueryRow(params)

	return result, nil
}

func QueryDB(query string, params ...string) (*sql.Rows, error) {
	stmt, err := Connection.Session.Prepare(query)
	if err != nil {
		return nil, err
	} else {
		defer stmt.Close()
	}

	results, err := stmt.Query(params)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func ExecDB(query string, params ...interface{}) (sql.Result, error) {
	stmt, err := Connection.Session.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(params...)
	if err != nil {
		return nil, err
	}
	stmt.Close()

	return result, nil
}
