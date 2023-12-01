package core

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBInfo struct {
	FilePath string
	// TODO
}

type dbSession struct {
	*sql.DB
}

func (info *DBInfo) buildConnectionString() string {
	connStr := info.FilePath
	// Configure the database connection string with the host, port, user, password, and dbname details
	return connStr
}

func openDBConnection(dbInfo DBInfo) dbSession {
	// Connect to postgres database and return session
	connectStr := dbInfo.buildConnectionString()
	fmt.Printf("Connect to sqlite database: %s\n", connectStr)
	db, err := sql.Open("sqlite3", connectStr)
	if err != nil {
		log.Panicf("Connect to database fail: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panicf("Cannot ping to database: %v", err)
	}

	// Optionally, you can use an ORM like GORM to simplify the database operations
	return dbSession{db}
}
