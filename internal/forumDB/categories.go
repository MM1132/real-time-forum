package forumDB

import (
	"database/sql"
)

type Category struct {
	CategoryID  int
	ParentID    sql.NullInt64
	Name        string
	Description sql.NullString

	ThreadCount int
	PostCount   int

	LatestPost Post
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

	categories := []Category{}
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
