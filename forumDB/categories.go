package forumDB

import (
	"database/sql"
	"forum/utils"
)

type Category struct {
	CategoryID  int
	ParentID    sql.NullInt64
	Name        string
	Description sql.NullString
}

type CategoryInterface interface {
	Insert(newCat Category) (int, error)
	Get(categoryID int) (Category, error)
	GetChildern(categoryID int) ([]Category, error)
}

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) Insert(newCat Category) (int, error) {
	stmt, err := m.DB.Prepare(
		"INSERT INTO categories(parentID, name, description) values(?,?,?)",
	)
	utils.FatalErr(err)

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
	stmt, err := m.DB.Prepare(
		"SELECT * FROM categories WHERE categoryID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(categoryID)
	category := Category{}
	err = row.Scan(
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
	stmt, err := m.DB.Prepare(
		"SELECT * FROM categories WHERE parentID=?",
	)
	utils.FatalErr(err)

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
