-- Func: Insert
INSERT INTO sessions(token, userID, created)
values (?, ?, ?);

-- Func: GetByToken
SELECT *
FROM sessions
WHERE token = ?;

-- Func: GetByUserID
SELECT *
FROM sessions
WHERE userID = ?;
