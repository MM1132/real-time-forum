package forumDB

import (
	"database/sql"
	"time"
)

type Board struct {
	BoardID     int
	ParentID    sql.NullInt64
	Name        string
	Description sql.NullString
	IsGroup     bool
	Order       int

	Extras *BoardExtras

	Children []Board
}

type BoardExtras struct {
	CountThreads int
	CountPosts   int

	LatestID       int
	LatestAuthorID int
	LatestAuthor   string
	LatestDate     time.Time

	ThreadID    int
	ThreadTitle string
}

type BoardModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewBoardModel(db *sql.DB) BoardModel {
	model := BoardModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/boards.sql")

	return model
}

func (m BoardModel) Insert(newBoard Board) (int, error) {
	stmt := m.statements["Insert"]

	res, err := stmt.Exec(
		newBoard.ParentID,
		newBoard.Name,
		newBoard.Description,
		newBoard.IsGroup,
		newBoard.Order,
	)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}

func (m BoardModel) Get(boardID int) (Board, error) {
	stmt := m.statements["Get"]

	row := stmt.QueryRow(boardID)
	board := Board{}
	err := row.Scan(
		&board.BoardID,
		&board.ParentID,
		&board.Name,
		&board.Description,
		&board.IsGroup,
		&board.Order,
	)
	if err != nil {
		return Board{}, err
	}

	return board, nil
}

func (m BoardModel) GetChildren(boardID int) ([]Board, error) {
	return m.getChildrenLocal(boardID, true, 0)
}

func (m BoardModel) getChildrenLocal(boardID int, extras bool, level int) ([]Board, error) {
	stmt := m.statements["GetChildren"]

	rows, err := stmt.Query(boardID)
	if err != nil {
		return nil, err
	}

	var boards []Board
	for rows.Next() {
		board := Board{}
		err = rows.Scan(
			&board.BoardID,
			&board.ParentID,
			&board.Name,
			&board.Description,
			&board.IsGroup,
			&board.Order,
		)
		if err != nil {
			return nil, err
		}

		if extras {
			board.Children, err = m.getChildrenLocal(board.BoardID, extras && board.IsGroup, level+1)
			if err != nil {
				return nil, err
			}

			if extras {
				board, err = m.GetExtras(board)
				if err != nil && err != sql.ErrNoRows {
					return nil, err
				}
			}
		}

		boards = append(boards, board)
	}

	return boards, nil
}

func (m BoardModel) GetBreadcrumbs(boardID int) ([]Board, error) {
	stmt := m.statements["GetBreadcrumbs"]

	rows, err := stmt.Query(boardID)
	if err != nil {
		return nil, err
	}

	var boards []Board
	for rows.Next() {
		board := Board{}
		err = rows.Scan(
			&board.BoardID,
			&board.ParentID,
			&board.Name,
			&board.Description,
			&board.IsGroup,
			&board.Order,
		)
		if err != nil {
			return nil, err
		}
		boards = append([]Board{board}, boards...)
	}
	return boards, nil
}

func (m BoardModel) SetSliceExtras(boards []Board) error {
	for i := range boards {
		if !boards[i].IsGroup {
			var err error
			boards[i], err = m.GetExtras(boards[i])
			if err != nil && err != sql.ErrNoRows {
				return err
			}
		}
	}

	return nil
}

func (m BoardModel) GetExtras(board Board) (Board, error) {
	stmt := m.statements["GetExtras"]

	row := stmt.QueryRow(board.BoardID)
	extras := BoardExtras{}
	err := row.Scan(
		&extras.CountThreads,
		&extras.CountPosts,

		&extras.LatestID,
		&extras.LatestAuthorID,
		&extras.LatestAuthor,
		&extras.LatestDate,

		&extras.ThreadID,
		&extras.ThreadTitle,
	)
	if err != nil {
		return board, err
	}

	board.Extras = &extras
	return board, nil
}
