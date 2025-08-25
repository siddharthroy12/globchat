package models

import "database/sql"

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
	stmt := "SELECT id, reason, reporter_id, message_id, created_at WHERE id = $1"

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
	return ReportQueryResult{}, nil
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
