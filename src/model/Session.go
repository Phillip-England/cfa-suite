package model

import (
	"cfa-suite/src/core"
	"crypto/rand"
	"encoding/base64"
)

type Session struct {
	ID int64
	UserID int64
	Token string
}

func NewSession() *Session {
	return &Session{}
}

func (m *Session) Insert(db *core.Database, userID int64) error {
	query := `
	INSERT INTO session (user_id, token)
	VALUES ($1, $2)
	RETURNING id
	`
	var sessionID int64
	randomBytes := make([]byte, 64)
	_, _ = rand.Read(randomBytes)
	token := base64.URLEncoding.EncodeToString(randomBytes)[:64]
	err := db.Connection.QueryRow(query, userID, token).Scan(&sessionID)
	if err != nil {
		return err
	}
	m.ID = sessionID
	m.UserID = userID
	m.Token = token
	return nil
}

func (m *Session) FindByToken(db *core.Database, token string) (error) {
	query := `
		SELECT * FROM session WHERE token = $1
	`
	err := db.Connection.QueryRow(query, token).Scan(&m.ID, &m.UserID, &m.Token)
	if err != nil {
		return err
	}
	return nil
}