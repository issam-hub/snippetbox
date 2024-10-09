package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	Expires   time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	statement := `INSERT INTO snippet (title, content, created_at, expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP, INTERVAL ? DAY))`

	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	statement := `SELECT * FROM snippet WHERE expires > UTC_TIMESTAMP() AND id = ?`

	s := &Snippet{}

	err := m.DB.QueryRow(statement, id).Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {

	statement := `SELECT * FROM snippet WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	snippets := []*Snippet{}

	rows, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.CreatedAt, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
