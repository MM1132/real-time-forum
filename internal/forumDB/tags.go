package forumDB

import (
	"database/sql"
	"fmt"
)

type TagModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewTagModel(db *sql.DB) TagModel {
	model := TagModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/tags.sql")

	return model
}

func (m TagModel) AddTags(threadID int, tags []string) error {
	if tags == nil || len(tags) == 0 {
		return nil
	}

	for _, tag := range tags {
		if err := m.Insert(threadID, tag); err != nil {
			return fmt.Errorf(`adding %v tag to %v: %w`, tag, threadID, err)
		}
	}
	return nil
}

func (m TagModel) Insert(threadID int, name string) error {
	stmt := m.statements["Insert"]
	_, err := stmt.Exec(name, threadID)
	return err
}

func (m TagModel) GetByThread(threadID int) (tags []string) {
	stmt := m.statements["GetByThread"]

	rows, err := stmt.Query(threadID)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return nil
		}

		tags = append(tags, tag)
	}

	return tags
}

func (m TagModel) GetThreads(tag string) (threads []Thread) {
	stmt := m.statements["GetThreads"]

	rows, err := stmt.Query(tag)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var thread Thread

		err = rows.Scan(
			&thread.ThreadID,
			&thread.Title,
			&thread.BoardID,
		)
		if err != nil {
			return nil
		}

		threads = append(threads, thread)
	}

	return threads
}

func (m TagModel) ThreadCount(tag string) int {
	stmt := m.statements["ThreadCount"]

	row := stmt.QueryRow(tag)

	var count int
	_ = row.Scan(&count)

	return count
}

func (m TagModel) GetPopular(boardID int) (tags []string) {
	stmt := m.statements["GetPopular"]

	rows, err := stmt.Query(boardID)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return nil
		}

		tags = append(tags, tag)
	}

	return tags
}
