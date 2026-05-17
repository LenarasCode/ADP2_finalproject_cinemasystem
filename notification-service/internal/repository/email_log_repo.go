package repository

import (
	"database/sql"
	"time"
)

type EmailLog struct {
	ID        string
	Recipient string
	Subject   string
	Body      string
	SentAt    time.Time
	Status    string
}

type EmailLogRepo struct {
	db *sql.DB
}

func NewEmailLogRepo(db *sql.DB) *EmailLogRepo {
	return &EmailLogRepo{db: db}
}

func (r *EmailLogRepo) Insert(log *EmailLog) error {
	_, err := r.db.Exec(
		`INSERT INTO email_logs (id, recipient, subject, body, sent_at, status) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		log.ID, log.Recipient, log.Subject, log.Body, log.SentAt, log.Status,
	)
	return err
}

func (r *EmailLogRepo) GetAll() ([]EmailLog, error) {
	rows, err := r.db.Query(`SELECT id, recipient, subject, body, sent_at, status FROM email_logs ORDER BY sent_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []EmailLog
	for rows.Next() {
		var l EmailLog
		if err := rows.Scan(&l.ID, &l.Recipient, &l.Subject, &l.Body, &l.SentAt, &l.Status); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}
