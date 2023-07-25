package model

import (
	"cfa-suite/src/core"
	"database/sql"
	"errors"
	"net/mail"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int64
	Email string
	Password string
}

func NewUser() (*User) {
	return &User{}
}

func (m *User) SetEmail(email string) (error) {
	m.Email = strings.ToLower(email)
	if m.Email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		return errors.New("invalid email address")
	}
	return nil
}

func (m *User) SetPassword(password string) (error) {
	if m.Password == "" {
		return errors.New("password is required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.Password = string(hashedPassword)
	return nil
}

func (m *User) SetID(id int64) (*User) {
	m.ID = id
	return m
}

func (m *User) IsUnique(database *core.Database) (bool, error) {
	query := `SELECT COUNT(*) FROM "user" WHERE email = $1`
	var count int
	err := database.Connection.QueryRow(query, m.Email).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	return true, nil
}

func (m *User) Insert(database *core.Database) error {
	statement := `INSERT INTO "user" (email, password) VALUES ($1, $2) RETURNING id`
	err := database.Connection.QueryRow(statement, m.Email, m.Password).Scan(&m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *User) FindByEmail(database *core.Database, email string) (*User, error) {
	query := `SELECT id, email, password FROM "user" WHERE email = $1`
	row := database.Connection.QueryRow(query, email)
	err := row.Scan(&m.ID, &m.Email, &m.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return m, nil
}

func (m *User) DeleteSessionsByUser(database *core.Database) error {
	query := `DELETE FROM session WHERE user_id = $1`
	_, err := database.Connection.Exec(query, m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}

func (m *User) Auth(c *gin.Context, database *core.Database) error {
	sessionToken, err := c.Cookie(os.Getenv("SESSION_TOKEN_KEY"))
	if err != nil {
		return err
	}
	sessionModel := NewSession()
	err = sessionModel.FindByToken(database, sessionToken)
	if err != nil {
		return err
	}
	m.SetID(sessionModel.UserID)
	trueUser := NewUser()
	err = trueUser.FindById(database, sessionModel.UserID)
	if err != nil {
		return err
	}
	if m.ID != trueUser.ID {
		return err
	}
	m.SetEmail(trueUser.Email)
	m.SetPassword(trueUser.Password)
	return nil
}

func (m *User) FindById(database *core.Database, userId int64) (error) {
	query := `SELECT id, email, password FROM "user" WHERE id = $1`
	row := database.Connection.QueryRow(query, userId)
	err := row.Scan(&m.ID, &m.Email, &m.Password)
	if err != nil {
		return err
	}
	return nil
}