package forumDB

import (
	"database/sql"
	"time"
)

type Category struct {
	CategoryID  int
	ParentID    sql.NullInt64
	Name        string
	Description sql.NullString

	Extras *CategoryExtras
}

type CategoryExtras struct {
	CountThreads int
	CountPosts   int

	LatestID       int
	LatestAuthorID int
	LatestAuthor   string
	LatestDate     time.Time

	ThreadID    int
	ThreadTitle string
}

type CategoryModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewCategoryModel(db *sql.DB) CategoryModel {
	model := CategoryModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/categories.sql")

	return model
}

func (m CategoryModel) Insert(newCat Category) (int, error) {
	stmt := m.statements["Insert"]

	res, err := stmt.Exec(
		newCat.ParentID,
		newCat.Name,
		newCat.Description,
	)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return int(id), nil
}

func (m CategoryModel) Get(categoryID int) (Category, error) {
	stmt := m.statements["Get"]

	row := stmt.QueryRow(categoryID)
	category := Category{}
	err := row.Scan(
		&category.CategoryID,
		&category.ParentID,
		&category.Name,
		&category.Description,
	)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (m CategoryModel) GetChildren(categoryID int) ([]Category, error) {
	stmt := m.statements["GetChildren"]

	rows, err := stmt.Query(categoryID)
	if err != nil {
		return nil, err
	}

	var categories []Category
	for rows.Next() {
		category := Category{}
		err = rows.Scan(
			&category.CategoryID,
			&category.ParentID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (m CategoryModel) GetBreadcrumbs(categoryID int) ([]Category, error) {
	stmt := m.statements["GetBreadcrumbs"]

	rows, err := stmt.Query(categoryID)
	if err != nil {
		return nil, err
	}

	var categories []Category
	for rows.Next() {
		category := Category{}
		err = rows.Scan(
			&category.CategoryID,
			&category.ParentID,
			&category.Name,
			&category.Description,
		)
		if err != nil {
			return nil, err
		}
		categories = append([]Category{category}, categories...)
	}
	return categories, nil
}

func (m CategoryModel) SetSliceExtras(categories []Category) error {
	for i := range categories {
		var err error
		categories[i], err = m.GetExtras(categories[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (m CategoryModel) GetExtras(category Category) (Category, error) {
	stmt := m.statements["GetExtras"]

	row := stmt.QueryRow(category.CategoryID)
	extras := CategoryExtras{}
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
		return Category{}, err
	}

	category.Extras = &extras
	return category, nil
}
