package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"
	"time"
)

type Thread struct {
	ID        int       `json:"id"`
	Lat       float64   `json:"lat"`
	Long      float64   `json:"long"`
	Message   string    `json:"message"`
	Replies   int       `json:"replies"`
	UserId    int       `json:"user_id"`
	Username  string    `json:"username"`
	UserImage string    `json:"user_image"`
	CreatedAt time.Time `json:"created_at"`
}

type ThreadModel struct {
	DB *sql.DB
}

func (m *ThreadModel) Create(message string, lat float64, long float64, userId int) (Thread, error) {
	stmt := "INSERT INTO threads (message, lat, long, user_id) VALUES($1, $2, $3, $4) RETURNING id, lat, long, message, user_id, created_at"

	var thread Thread
	err := m.DB.QueryRow(stmt, message, lat, long, userId).Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message, &thread.UserId, &thread.CreatedAt)

	if err != nil {
		return Thread{}, err
	}

	var user User

	stmt = "SELECT username, image FROM users WHERE id = $1"
	err = m.DB.QueryRow(stmt, userId).Scan(&user.Username, &user.Image)

	if err != nil {
		return Thread{}, err
	}
	thread.Username = user.Username
	thread.UserImage = user.Image
	return thread, nil
}

func (m *ThreadModel) Delete(threadId int) error {
	stmt := "DELETE FROM threads WHERE id = $1"

	result, err := m.DB.Exec(stmt, threadId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *ThreadModel) GetById(threadId int) (Thread, error) {
	stmt := `SELECT threads.id, lat, long, message, user_id, threads.created_at, 
             users.username, users.image 
             FROM threads 
             INNER JOIN users ON users.id = threads.user_id 
             WHERE threads.id = $1`

	row := m.DB.QueryRow(stmt, threadId)

	thread := Thread{}
	err := row.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message,
		&thread.UserId, &thread.CreatedAt, &thread.Username, &thread.UserImage)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Thread{}, ErrNoRecord // Convert to your custom error
		}
		return Thread{}, err
	}

	return thread, nil
}

func (m *ThreadModel) GetRandomThread() (Thread, error) {
	// Get total count first
	var count int
	countStmt := "SELECT COUNT(*) FROM threads"
	err := m.DB.QueryRow(countStmt).Scan(&count)
	if err != nil {
		return Thread{}, err
	}

	if count == 0 {
		return Thread{}, sql.ErrNoRows
	}

	// Generate cryptographically secure random offset
	maxBig := big.NewInt(int64(count))
	randomOffset, err := rand.Int(rand.Reader, maxBig)
	if err != nil {
		return Thread{}, err
	}

	stmt := `SELECT threads.id, lat, long, message, user_id, threads.created_at, users.username, users.image 
			 FROM threads 
			 INNER JOIN users ON users.id = threads.user_id 
			 LIMIT 1 OFFSET $1`

	row := m.DB.QueryRow(stmt, randomOffset.Int64())

	thread := Thread{}
	err = row.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message, &thread.UserId, &thread.CreatedAt, &thread.Username, &thread.UserImage)

	if err != nil {
		return Thread{}, err
	}
	return thread, nil
}

func (m *ThreadModel) IncreaseReplies(threadId int) error {
	stmt := "UPDATE threads SET replies = replies + 1 WHERE id = $1"

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

func (m *ThreadModel) DecreaseReplies(threadId int) error {
	stmt := "UPDATE threads SET replies = replies - 1 WHERE id = $1"

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

func (m *ThreadModel) SetReplies(threadId int, replies int) error {
	stmt := "UPDATE threads SET replies = $1 WHERE id = $2"

	result, err := m.DB.Exec(stmt, replies, threadId)
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

func (m *ThreadModel) GetAllByUserId(userId int) ([]*Thread, error) {
	stmt := "SELECT threads.id, lat, long, message, user_id, threads.created_at, users.username, users.image FROM threads INNER JOIN users ON users.id = threads.user_id WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []*Thread{}

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message, &thread.UserId, &thread.CreatedAt, &thread.Username, &thread.UserImage)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (m *ThreadModel) GetByLocation(minLat, maxLat, minLong, maxLong float64) ([]*Thread, error) {
	stmt := `SELECT threads.id, lat, long, message, user_id, threads.created_at, users.username, users.image
			 FROM threads INNER JOIN users ON users.id = threads.user_id
			 WHERE lat BETWEEN $1 AND $2 
			 AND long BETWEEN $3 AND $4 
			 ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, minLat, maxLat, minLong, maxLong)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []*Thread{}

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message, &thread.UserId, &thread.CreatedAt)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (m *ThreadModel) GetByLocationRadius(centerLat, centerLong, radiusKm float64) ([]*Thread, error) {
	// Using the spherical law of cosines for distance calculation
	// This is an approximation suitable for most use cases
	stmt := `SELECT threads.id, lat, long, message, user_id, threads.created_at, users.username, users.image
			 FROM threads INNER JOIN users ON users.id = threads.user_id
			 WHERE (
				6371 * acos(
					cos(radians($1)) * cos(radians(lat)) * 
					cos(radians(long) - radians($2)) + 
					sin(radians($1)) * sin(radians(lat))
				)
			 ) <= $3
			 ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, centerLat, centerLong, radiusKm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := make([]*Thread, 0)

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.Message, &thread.UserId, &thread.CreatedAt, &thread.Username, &thread.UserImage)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return threads, nil
}
