package models

import (
	"database/sql"
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

func (m *MessageModel) GetByID(messageId int) (Message, error) {
	stmt := "SELECT messages.id, text, messages.image, thread_id, reported, is_first, user_id, messages.created_at, users.username, users.image FROM messages INNER JOIN users ON users.id = messages.user_id WHERE messages.id = $1"

	message := Message{}
	err := m.DB.QueryRow(stmt, messageId).Scan(&message.ID, &message.Text, &message.Image, &message.ThreadId, &message.Reported, &message.IsFirst, &message.UserId, &message.CreatedAt, &message.Username, &message.UserImage)

	if err != nil {
		return Message{}, err
	}

	return message, nil
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
