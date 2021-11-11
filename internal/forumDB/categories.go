package forumDB

import (
	"database/sql"
	"forum/internal/utils"
)

type Category struct {
	CategoryID  int
	ParentID    sql.NullInt64
	Name        string
	Description sql.NullString
}

type CategoryModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewCategoryModel(db *sql.DB) CategoryModel {
	statements := make(map[string]*sql.Stmt)
	model := CategoryModel{db: db}

	var err error
	statements["Insert"], err = db.Prepare("INSERT INTO categories(parentID, name, description) values(?,?,?)")
	utils.FatalErr(err)

	statements["Get"], err = db.Prepare("SELECT * FROM categories WHERE categoryID=?")
	utils.FatalErr(err)

	statements["GetChildern"], err = db.Prepare("SELECT * FROM categories WHERE parentID=?")
	utils.FatalErr(err)

	statements["GetBreadcrumbs"], err = db.Prepare(`
		WITH ancestors AS (
			SELECT *
			FROM categories
			WHERE categoryID=?
			
			UNION ALL
		
			SELECT c.*
			FROM categories c
				JOIN
				ancestors a ON c.categoryID = a.parentID
		)
		SELECT * FROM ancestors`)
	utils.FatalErr(err)

	model.statements = statements
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

func (m CategoryModel) GetChildern(categoryID int) ([]Category, error) {
	stmt := m.statements["GetChildern"]

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
