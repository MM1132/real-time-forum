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

-- Func: GetByEmail
SELECT *
FROM users
WHERE email = ?;

-- Func: SetExtras
SELECT count(*)
FROM posts
WHERE userID = ?;

-- Func: SetImage
UPDATE users
SET image = ?
WHERE userID = ?;