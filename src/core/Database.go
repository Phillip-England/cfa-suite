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
	err = d.CreateEmailKeyTable()
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
		password TEXT NOT NULL,
		active BOOLEAN NOT NULL
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
			number VARCHAR(255),
			FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
		)
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateEmailKeyTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS email_key (
			id SERIAL PRIMARY KEY,
			user_id INT,
			key VARCHAR(255),
			FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
		)
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddActiveColumnToUserTable() error {
	query := `
		ALTER TABLE "user" ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT false
	`
	_, err := d.Connection.Exec(query)
	if err != nil {
		return err
	}

	// Set all users' 'active' field to 'false'
	_, err = d.Connection.Exec("UPDATE \"user\" SET active = false")
	if err != nil {
		return err
	}

	return nil
}