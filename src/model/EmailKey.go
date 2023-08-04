package model

import (
	"cfa-suite/src/core"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailKey struct {
	ID int64
	UserID *int64
	Key string
}

func NewEmailKey() *EmailKey {
	return &EmailKey{}
}

func (m *EmailKey) SetUserID(userID int64) {
	m.UserID = &userID
}

func (m *EmailKey) GetUserID() (*int64, error) {
	if m.UserID == nil {
		return nil, errors.New("userID not set in EmailKey")
	}
	return m.UserID, nil
}

func (m *EmailKey) GenerateRandomKey() {
	randomBytes := make([]byte, 64)
	_, _ = rand.Read(randomBytes)
	key := base64.URLEncoding.EncodeToString(randomBytes)[:64]
	m.Key = key
}

func (m *EmailKey) Insert(database *core.Database) (error) {
	m.GenerateRandomKey()
	statement := `INSERT INTO email_key (user_id, key) VALUES ($1, $2) RETURNING id`
	err := database.Connection.QueryRow(statement, m.UserID, m.Key).Scan(&m.ID)
	if err != nil {
		return err
	}
	return nil
}

func (m *EmailKey) SendAccountVerificationEmail(email string) error {
    if m.Key == "" {
        return errors.New("verification key not generated")
    }
	from := mail.NewEmail("Example User", os.Getenv("APP_EMAIL"))
	subject := "Email Verification |"
	to := mail.NewEmail("CFA Suite", "phillip.mark.england@gmail.com")
	var verificationLink string
	if os.Getenv("GO_ENV") == "dev" {
		verificationLink = "http://" + os.Getenv("SERVER_URL") + ":" + os.Getenv("PORT") + "/api/verify-account/" + m.Key
	} else {
		verificationLink = "https://" + os.Getenv("SERVER_URL") + "/api/verify-account/" + m.Key
	}
	body := fmt.Sprintf("Dear user,\n\nPlease use the following verification link to verify your account:\n\n%s\n\nBest regards,\nThe CFA Suite Team", verificationLink)
	message := mail.NewSingleEmail(from, subject, to, body, "")
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (m *EmailKey) FindByKey(database *core.Database, key string) (error, bool) {
	statement := `SELECT id, user_id, key FROM email_key WHERE key = $1 LIMIT 1`
	row := database.Connection.QueryRow(statement, key)
	var id int64
	var userID *int64
	var dbKey string
	err := row.Scan(&id, &userID, &dbKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false
		}
		return err, false
	}
	m.ID = id
	m.UserID = userID
	m.Key = dbKey
	return nil, true
}