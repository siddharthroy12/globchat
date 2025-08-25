package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
	Username  string    `json:"username"`
	Image     string    `json:"image"`
	Messages  int       `json:"messages"`
}

type UserQuery struct {
	Search    string
	PageSize  int
	PageIndex int
}

type UserQueryResult struct {
	Total int    `json:"total"`
	Count int    `json:"count"`
	Users []User `json:"users"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(email string, username string) (User, error) {
	stmt := "INSERT INTO users (email, username) VALUES($1, $2) RETURNING id, email, created_at, username, image, messages, is_admin"

	row := m.DB.QueryRow(stmt, email, username)

	user, err := m.getUserFromRow(row)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m *UserModel) getUserFromRow(row *sql.Row) (User, error) {
	var u User

	err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.Username, &u.Image, &u.Messages, &u.IsAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}

	return u, nil
}

func (m *UserModel) GetById(userId int) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages, is_admin FROM users WHERE id = $1"
	row := m.DB.QueryRow(stmt, userId)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetByEmail(email string) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages, is_admin FROM users WHERE email = $1"
	row := m.DB.QueryRow(stmt, email)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetByUsername(username string) (User, error) {
	stmt := "SELECT id, email, created_at, username, image, messages, is_admin FROM users WHERE username = $1"
	row := m.DB.QueryRow(stmt, username)
	return m.getUserFromRow(row)
}

func (m *UserModel) GetFromSessionToken(token string) (User, error) {
	stmt := `SELECT users.id, users.email, users.created_at, users.username, users.image, users.messages , is_admin
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

func (m *UserModel) Query(query UserQuery) (UserQueryResult, error) {
	// Build the base query
	baseStmt := `SELECT id, email, created_at, username, image, messages, is_admin FROM users`
	countStmt := `SELECT COUNT(*) FROM users`

	var whereClause string
	var args []interface{}

	// Add search condition if provided
	if query.Search != "" {
		whereClause = ` WHERE username ILIKE $1 OR email ILIKE $1`
		args = append(args, "%"+query.Search+"%")
	}

	// Get total count
	var total int
	countQuery := countStmt + whereClause
	err := m.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return UserQueryResult{}, err
	}

	// Build the main query with pagination
	mainQuery := baseStmt + whereClause + ` ORDER BY created_at DESC`

	// Add pagination
	if query.PageSize > 0 {
		offset := query.PageIndex * query.PageSize
		mainQuery += ` LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
		args = append(args, query.PageSize, offset)
	}

	// Execute the query
	rows, err := m.DB.Query(mainQuery, args...)
	if err != nil {
		return UserQueryResult{}, err
	}
	defer rows.Close()

	// Scan results
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.Username, &u.Image, &u.Messages, &u.IsAdmin)
		if err != nil {
			return UserQueryResult{}, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return UserQueryResult{}, err
	}

	result := UserQueryResult{
		Total: total,
		Count: len(users),
		Users: users,
	}

	return result, nil
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
