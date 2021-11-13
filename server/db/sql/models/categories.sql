-- Func: Insert
INSERT INTO categories(parentID, name, description) values(?,?,?);

-- Func: Get
SELECT * FROM categories WHERE categoryID=?;

-- Func: GetChildren
SELECT * FROM categories WHERE parentID=?;

-- Func: GetBreadcrumbs
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
SELECT * FROM ancestors;