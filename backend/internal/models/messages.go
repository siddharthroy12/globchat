package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Image     string    `json:"image"`
	ThreadId  int       `json:"thread_id"`
	Reported  int       `json:"reported"`
	IsFirst   bool      `json:"is_first"`
	UserId    int       `json:"user_id"`
	Username  string    `json:"username"`
	UserImage string    `json:"user_image"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageQuery struct {
	Search    string
	PageSize  int
	PageIndex int
}

type MessageQueryResult struct {
	Total    int       `json:"total"`
	Count    int       `json:"count"`
	Messages []Message `json:"messages"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Create(text string, image string, threadId int, userId int, isFirst bool) (Message, error) {
	stmt := "INSERT INTO messages (text, image, thread_id, user_id, is_first) VALUES($1, $2, $3, $4, $5) RETURNING id, text, image, thread_id, reported, is_first, user_id, created_at"

	var message Message
	err := m.DB.QueryRow(stmt, text, image, threadId, userId, isFirst).Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt)

	if err != nil {
		return Message{}, err
	}

	var user User

	stmt = "SELECT username, image FROM users WHERE id = $1"
	err = m.DB.QueryRow(stmt, userId).Scan(&user.Username, &user.Image)

	if err != nil {
		return Message{}, err
	}
	message.Username = user.Username
	message.UserImage = user.Image

	return message, nil
}

func (m *MessageModel) Delete(messageId int) error {
	stmt := "DELETE FROM messages WHERE id = $1"

	result, err := m.DB.Exec(stmt, messageId)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the message didn't exist
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *MessageModel) DeleteByThreadID(threadId int) error {
	stmt := "DELETE FROM messages WHERE thread_id = $1"

	result, err := m.DB.Exec(stmt, threadId)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the message didn't exist
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *MessageModel) IncreaseReported(messageId int) error {
	stmt := "UPDATE messages SET reported = reported + 1 WHERE id = $1"

	result, err := m.DB.Exec(stmt, messageId)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were affected, the message didn't exist
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (m *MessageModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM messages WHERE id = $1)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (m *MessageModel) GetByID(messageId int) (Message, error) {
	stmt := "SELECT messages.id, text, messages.image, thread_id, reported, is_first, user_id, messages.created_at, users.username, users.image FROM messages INNER JOIN users ON users.id = messages.user_id WHERE messages.id = $1"

	message := Message{}
	err := m.DB.QueryRow(stmt, messageId).Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt, &message.Username, &message.UserImage)

	if err != nil {
		return Message{}, err
	}

	return message, nil
}

func (m *MessageModel) Query(query MessageQuery) (MessageQueryResult, error) {
	var result MessageQueryResult

	// Base query with JOIN to get user information
	baseQuery := `SELECT messages.id, text, messages.image, thread_id, reported, is_first, 
	              user_id, messages.created_at, users.username, users.image 
	              FROM messages 
	              INNER JOIN users ON users.id = messages.user_id`

	// Count query for total results
	countQuery := `SELECT COUNT(*) FROM messages 
	               INNER JOIN users ON users.id = messages.user_id`

	var whereClause string
	var args []interface{}
	argIndex := 1

	// Add search condition if provided
	if query.Search != "" {
		whereClause = " WHERE LOWER(text) LIKE LOWER($" + fmt.Sprintf("%d", argIndex) + ")"
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	// Get total count first
	countStmt := countQuery + whereClause
	err := m.DB.QueryRow(countStmt, args...).Scan(&result.Total)
	if err != nil {
		return MessageQueryResult{}, err
	}

	// Calculate offset for pagination
	offset := query.PageIndex * query.PageSize

	// Build final query with pagination
	finalQuery := baseQuery + whereClause + " ORDER BY messages.created_at DESC LIMIT $" +
		fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)

	args = append(args, query.PageSize, offset)

	// Execute the query
	rows, err := m.DB.Query(finalQuery, args...)
	if err != nil {
		return MessageQueryResult{}, err
	}
	defer rows.Close()

	var messages []Message

	// Scan results
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId,
			&message.Reported, &message.IsFirst, &message.UserId,
			&message.CreatedAt, &message.Username, &message.UserImage)
		if err != nil {
			return MessageQueryResult{}, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return MessageQueryResult{}, err
	}

	result.Messages = messages
	result.Count = len(messages)

	return result, nil
}

func (m *MessageModel) GetByThreadID(threadId int, limit int) ([]Message, error) {
	stmt := "SELECT messages.id, text, messages.image, thread_id, reported, is_first, user_id, messages.created_at, users.username, users.image FROM messages INNER JOIN users ON users.id = messages.user_id WHERE thread_id = $1 ORDER BY created_at DESC LIMIT $2"

	rows, err := m.DB.Query(stmt, threadId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt, &message.Username, &message.UserImage)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageModel) GetBeforeID(threadId int, id int, limit int) ([]Message, error) {
	stmt := "SELECT messages.id, text, messages.image, thread_id, reported, is_first, user_id, messages.created_at, users.username, users.image FROM messages INNER JOIN users ON users.id = messages.user_id WHERE thread_id = $1 AND messages.id < $2 ORDER BY messages.id DESC LIMIT $3"

	rows, err := m.DB.Query(stmt, threadId, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt, &message.Username, &message.UserImage)

		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *MessageModel) GetAfterID(threadId int, id int, limit int) ([]Message, error) {
	stmt := "SELECT messages.id, text, messages.image, thread_id, reported, is_first, user_id, messages.created_at, users.username, users.image FROM messages INNER JOIN users ON users.id = messages.user_id WHERE thread_id = $1 AND messages.id > $2 ORDER BY messages.id ASC LIMIT $3"

	rows, err := m.DB.Query(stmt, threadId, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}

	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt, &message.Username, &message.UserImage)

		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
