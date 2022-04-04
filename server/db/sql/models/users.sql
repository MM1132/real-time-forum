-- Func: Insert
INSERT INTO users(nickname, firstname, lastname, age, gender, email, password, created)
values (?, ?, ?, ?, ?, ?, ?, ?);

-- Func: Get
SELECT *
FROM users
WHERE userID = ?;

-- Func: GetByName
SELECT *
FROM users
WHERE nickname = ?;

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

-- Func: SetPassword
UPDATE users
SET password = ?
WHERE userID = ?;

-- Func: SetDescription
UPDATE users
SET description = ?
WHERE userID = ?;
