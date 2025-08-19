package models

import (
	"database/sql"
	"time"

	"globechat.live/internal/crypto"
)

type Session struct {
	UserId    int
	CreatedAt time.Time
	ExpiresAt time.Time
	Token     string
}

func (s Session) HasExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

type SessionModel struct {
	DB *sql.DB
}

func (m *SessionModel) Create(
	userId int,
) (string, error) {

	err := m.Remove(userId)
	if err != nil {
		return "", err
	}

	token, err := crypto.GenerateRandomToken(32)
	if err != nil {
		return "", err
	}

	stmt := "INSERT INTO sessions (user_id, token) VALUES($1, $2)"

	_, err = m.DB.Exec(stmt, userId, token)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *SessionModel) Exists(userId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE user_id = $1)"

	err := m.DB.QueryRow(stmt, userId).Scan(&exists)

	return exists, err
}

func (m *SessionModel) Remove(
	userId int,
) error {

	stmt := "DELETE FROM sessions WHERE user_id = $1"

	_, err := m.DB.Exec(stmt, userId)

	if err != nil {
		return err
	}

	return nil
}
