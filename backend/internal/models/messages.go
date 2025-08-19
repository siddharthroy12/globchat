package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	ImageUrl  string    `json:"image_url"`
	ThreadId  int       `json:"thread_id"`
	Reported  int       `json:"reported"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Create(text string, imageUrl string, threadId int, userId int) (int, error) {
	stmt := "INSERT INTO messages (text, image_url, thread_id, user_id) VALUES($1, $2, $3, $4) RETURNING id"

	var id int
	err := m.DB.QueryRow(stmt, text, imageUrl, threadId, userId).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
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

func (m *MessageModel) GetByID(messageId int) (*Message, error) {
	stmt := "SELECT id, text, image_url, thread_id, user_id, created_at FROM messages WHERE id = $1"

	message := &Message{}
	err := m.DB.QueryRow(stmt, messageId).Scan(
		&message.ID,
		&message.Text,
		&message.ImageUrl,
		&message.ThreadId,
		&message.UserId,
		&message.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func (m *MessageModel) GetByThreadID(threadId int, limit int) ([]*Message, error) {
	stmt := "SELECT id, text, image_url, thread_id, user_id, created_at FROM messages WHERE thread_id = $1 ORDER BY created_at DESC LIMIT $2"

	rows, err := m.DB.Query(stmt, threadId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		message := &Message{}
		err = rows.Scan(
			&message.ID,
			&message.Text,
			&message.ImageUrl,
			&message.ThreadId,
			&message.UserId,
			&message.CreatedAt,
		)
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

// GetBefore gets n messages before time t in a thread
func (m *MessageModel) GetBeforeID(threadId int, id int, limit int) ([]*Message, error) {
	stmt := "SELECT id, text, image_url, thread_id, user_id, created_at FROM messages WHERE thread_id = $1 AND id < $2 ORDER BY id DESC LIMIT $3"

	rows, err := m.DB.Query(stmt, threadId, id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}

	for rows.Next() {
		message := &Message{}
		err = rows.Scan(
			&message.ID,
			&message.Text,
			&message.ImageUrl,
			&message.ThreadId,
			&message.UserId,
			&message.CreatedAt,
		)
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
