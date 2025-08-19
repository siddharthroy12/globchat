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
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(email string, username string) (int, error) {
	stmt := "INSERT INTO users (email, username) VALUES($1, $2) RETURNING id"

	var id int
	result, err := m.DB.Query(stmt, email, username)

	if err != nil {
		return 0, err
	}

	err = result.Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) getUserFromRow(row *sql.Row) (User, error) {
	var u User

	err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}

	return u, nil
}

func (m *UserModel) GetById(userId int) (User, error) {
	stmt := "SELECT id, email, created_at, username FROM users WHERE id = $1"

	row := m.DB.QueryRow(stmt, userId)

	return m.getUserFromRow(row)
}

func (m *UserModel) GetByEmail(email string) (User, error) {
	stmt := "SELECT id, email, created_at, username FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)

	return m.getUserFromRow(row)
}

func (m *UserModel) GetByUsername(username string) (User, error) {
	stmt := "SELECT id, email, created_at, username FROM users WHERE username = $1"
	row := m.DB.QueryRow(stmt, username)

	return m.getUserFromRow(row)
}

func (m *UserModel) GetFromSessionToken(token string) (User, error) {
	stmt := `SELECT users.id, users.email, users.created_at, users.username FROM sessions INNER JOIN users ON users.id = sessions.user_id WHERE sessions.token = $1`
	row := m.DB.QueryRow(stmt, token)

	return m.getUserFromRow(row)
}

func (m *UserModel) Delete(userId int) error {
	stmt := "DELETE FROM users WHERE id = $1"

	_, err := m.DB.Exec(stmt, userId)

	if err != nil {
		return err
	}

	return nil
}
