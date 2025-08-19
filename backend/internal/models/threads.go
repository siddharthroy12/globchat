package models

import (
	"database/sql"
	"time"
)

type Thread struct {
	ID        int       `json:"id"`
	Lat       float64   `json:"lat"`
	Long      float64   `json:"long"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ThreadModel struct {
	DB *sql.DB
}

func (m *ThreadModel) Create(lat float64, long float64, userId int) (int, error) {
	stmt := "INSERT INTO threads (lat, long, user_id) VALUES($1, $2, $3) RETURNING id"

	var id int
	err := m.DB.QueryRow(stmt, lat, long, userId).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ThreadModel) Delete(threadId int) error {
	stmt := "DELETE FROM threads WHERE id = $1"

	_, err := m.DB.Exec(stmt, threadId)

	if err != nil {
		return err
	}

	return nil
}

func (m *ThreadModel) GetAllByUserId(userId int) ([]*Thread, error) {
	stmt := "SELECT id, lat, long, user_id, created_at FROM threads WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threads := []*Thread{}

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.UserId, &thread.CreatedAt)
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
	stmt := `SELECT id, lat, long, user_id, created_at 
			 FROM threads 
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
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.UserId, &thread.CreatedAt)
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
	stmt := `SELECT id, lat, long, user_id, created_at 
			 FROM threads 
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

	threads := []*Thread{}

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.ID, &thread.Lat, &thread.Long, &thread.UserId, &thread.CreatedAt)
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
