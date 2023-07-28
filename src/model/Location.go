package model

import (
	"cfa-suite/src/core"
	"errors"
	"strconv"
)

type Location struct {
	ID int64
	UserID int64
	Name string
	Number string
}

func NewLocation() *Location {
	return &Location{}
}

func (m *Location) SetUserID(userID int64) {
	m.UserID = userID
}

func (m *Location) SetName(name string) error {
	if len(name) > 64 {
		return errors.New("name too long")
	}
	if len(name) < 3 {
		return errors.New("name too short")
	}
	m.Name = name
	return nil
}

func (m *Location) SetNumber(stringNumber string) error {
	if len(stringNumber) > 12 {
		return errors.New("number too long")
	}
	if len(stringNumber) < 3 {
		return errors.New("number too short")
	}
	_, err := strconv.ParseInt(stringNumber, 10, 64)
	if err != nil {
		return errors.New("please provide a valid number")
	}
	m.Number = stringNumber
	return nil
}

func (m *Location) Insert(database *core.Database) error {
	statement := `INSERT INTO location (user_id, name, number) VALUES ($1, $2, $3) RETURNING id`
	err := database.Connection.QueryRow(statement, m.UserID, m.Name, m.Number).Scan(&m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Location) GetByID(locationID string, database *core.Database) error {
	query := `SELECT id, user_id, name, number FROM location WHERE id = $1`
	err := database.Connection.QueryRow(query, locationID).Scan(&m.ID, &m.UserID, &m.Name, &m.Number)
	if err != nil {
		return err
	}
	return nil
}

func (m *Location) Update(database *core.Database) error {
	statement := `UPDATE location SET name = $1, number = $2 WHERE id = $3`
	_, err := database.Connection.Exec(statement, m.Name, m.Number, m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Location) Delete(database *core.Database) error {
	statement := `DELETE FROM location WHERE id = $1`
	_, err := database.Connection.Exec(statement, m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *Location) LimitNumberOfLocations(userID int64, database *core.Database) (bool, error) {
	query := `SELECT COUNT(*) FROM location WHERE user_id = $1`
	var count int
	err := database.Connection.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count >= 3, nil
}

func (m *Location) GetLocationsByUserID(userID int64, database *core.Database) ([]*Location, error) {
	query := `SELECT id, user_id, name, number FROM location WHERE user_id = $1`
	rows, err := database.Connection.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := []*Location{}
	for rows.Next() {
		location := NewLocation()
		err := rows.Scan(&location.ID, &location.UserID, &location.Name, &location.Number)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}