package model

import "errors"

type EmailKey struct {
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