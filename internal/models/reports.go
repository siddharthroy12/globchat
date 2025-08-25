package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type Report struct {
	ID         int    `json:"id"`
	Reason     string `json:"reason"`
	ReporterId int    `json:"reporter_id"`
	MessageId  int    `json:"message_id"`
	CreatedAt  string `json:"created_at"`
}

type ReportQuery struct {
	Search    string
	PageSize  int
	PageIndex int
}

type ReportQueryResult struct {
	Total   int      `json:"total"`
	Count   int      `json:"count"`
	Reports []Report `json:"reports"`
}

type ReportModel struct {
	DB *sql.DB
}

func (m *ReportModel) Create(
	userId int, messageId int, reason string,
) error {

	stmt := "INSERT INTO reports (reporter_id, message_id, reason) VALUES($1, $2, $3)"

	_, err := m.DB.Exec(stmt, userId, messageId, reason)

	if err != nil {
		return err
	}

	return nil
}

func (m *ReportModel) Exists(userId int, messageId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM reports WHERE reporter_id = $1 AND message_id = $2)"

	err := m.DB.QueryRow(stmt, userId, messageId).Scan(&exists)

	return exists, err
}

func (m *ReportModel) GetByID(reportId int) (Report, error) {
	stmt := "SELECT id, reason, reporter_id, message_id, created_at FROM reports WHERE id = $1"

	report := Report{}
	err := m.DB.QueryRow(stmt, reportId).Scan(&report.ID, &report.Reason, &report.ReporterId, &report.MessageId, &report.CreatedAt)

	if err != nil {
		return Report{}, err
	}

	return report, nil
}

func (m *ReportModel) RemoveById(
	id int,
) error {

	stmt := "DELETE FROM reports WHERE id = $1"

	_, err := m.DB.Exec(stmt, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *ReportModel) RemoveByReporterId(
	userId int,
) error {

	stmt := "DELETE FROM reports WHERE reporter_id = $1"

	_, err := m.DB.Exec(stmt, userId)

	if err != nil {
		return err
	}

	return nil
}

func (m *ReportModel) Query(query ReportQuery) (ReportQueryResult, error) {
	var result ReportQueryResult
	var args []interface{}
	var conditions []string
	argIndex := 1

	// Build the WHERE clause for search
	baseQuery := "FROM reports"
	if query.Search != "" {
		conditions = append(conditions, fmt.Sprintf("reason ILIKE $%d", argIndex))
		args = append(args, "%"+query.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	err := m.DB.QueryRow(countQuery, args...).Scan(&result.Total)
	if err != nil {
		return ReportQueryResult{}, err
	}

	// Build the main query with pagination
	selectQuery := "SELECT id, reason, reporter_id, message_id, created_at " + baseQuery + whereClause + " ORDER BY created_at DESC"

	// Add pagination
	if query.PageSize > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, query.PageSize)
		argIndex++

		if query.PageIndex > 0 {
			selectQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, query.PageIndex*query.PageSize)
		}
	}

	// Execute the main query
	rows, err := m.DB.Query(selectQuery, args...)
	if err != nil {
		return ReportQueryResult{}, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var report Report
		err := rows.Scan(&report.ID, &report.Reason, &report.ReporterId, &report.MessageId, &report.CreatedAt)
		if err != nil {
			return ReportQueryResult{}, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return ReportQueryResult{}, err
	}

	result.Reports = reports
	result.Count = len(reports)

	return result, nil
}

func (m *ReportModel) RemoveByMessageId(
	messageId int,
) error {

	stmt := "DELETE FROM reports WHERE message_id = $1"

	_, err := m.DB.Exec(stmt, messageId)

	if err != nil {
		return err
	}

	return nil
}
