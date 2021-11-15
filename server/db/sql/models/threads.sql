-- Func: Insert
INSERT INTO threads(title, categoryID)
values (?, ?);

-- Func: Get
SELECT *
FROM threads
WHERE threadID = ?;

-- Func: ByCategory
SELECT *
FROM threads
WHERE categoryID = ?;
