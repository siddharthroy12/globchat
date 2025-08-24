package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int
	Email     string
	CreatedAt time.Time
	Username  string
	Image     string
	Messages  int
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(email string, username string) (User, error) {
	stmt := "INSERT INTO users (email, username) VALUES($1, $2) RETURNING id, email, created_at, username, image, messages"

	row := m.DB.QueryRow(stmt, email, username)

	user, err := m.getUserFromRow(row)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *UserModel) getUserFromRow(row *sql.Row) (User, error) {
	var u User

	err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.Username, &u.Image, &u.Messages)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}

	return u, nil
}

func (m *UserModel) GetById(userId int) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages FROM users WHERE id = $1"
	row := m.DB.QueryRow(stmt, userId)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetByEmail(email string) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetByUsername(username string) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages FROM users WHERE username = $1"
	row := m.DB.QueryRow(stmt, username)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetFromSessionToken(token string) (User, error) {
	stmt := `SELECT users.id, users.email, users.created_at, users.username, users.image, users.messages 
	         FROM sessions 
	         INNER JOIN users ON users.id = sessions.user_id 
	         WHERE sessions.token = $1 AND sessions.expires_at > NOW()`
	row := m.DB.QueryRow(stmt, token)
	return m.getUserFromRow(row)
}

func (m *UserModel) UpdateImageAndUsername(userId int, image string, username string) error {
	stmt := "UPDATE users SET image = $1, username = $2 WHERE id = $3"

	_, err := m.DB.Exec(stmt, image, username, userId)

	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) UpdateImage(userId int, image string) error {
	stmt := "UPDATE users SET image = $1 WHERE id = $2"
	_, err := m.DB.Exec(stmt, image, userId)
	return err
}

func (m *UserModel) UpdateMessages(userId int, messages int) error {
	stmt := "UPDATE users SET messages = $1 WHERE id = $2"
	_, err := m.DB.Exec(stmt, messages, userId)
	return err
}

func (m *UserModel) Delete(userId int) error {
	stmt := "DELETE FROM users WHERE id = $1"

	_, err := m.DB.Exec(stmt, userId)

	if err != nil {
		return err
	}

	return nil
}
