package core

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
    Connection *sql.DB
}

func NewDatabase() (*Database) {
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
    if err != nil {
        log.Panic(err.Error())
    }
    err = db.Ping()
    if err != nil {
        log.Panic(err.Error())
    }
    return &Database{
        Connection: db,
    }
}

func (d *Database) InitTables() error {
	err := d.CreateUserTable()
    if err != nil {
        return err
    }
    err = d.CreateSessionTable()
    if err != nil {
		return err
    }
	err = d.CreateLocationTable()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
		)
	`
	_, err := d.Connection.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

func (d *Database) CreateSessionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS session (
			id SERIAL PRIMARY KEY,
			user_id INT, -- Add the user_id column,
            token TEXT,
			FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
		)
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateLocationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS location (
			id SERIAL PRIMARY KEY,
			user_id INT,
			name VARCHAR(255),
			number INT,
			FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
		)
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddColumnsToUserTable() error {
	query := `
		ALTER TABLE "user"
		DROP COLUMN first_name,
		DROP COLUMN last_name
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
