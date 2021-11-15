-- Func: Insert
INSERT INTO users(name, email, password, created)
values (?, ?, ?, ?);

-- Func: Get
SELECT *
FROM users
WHERE userID = ?;

-- Func: GetByName
SELECT *
FROM users
WHERE name = ?;
