-- Func: Insert
INSERT INTO threads(title, boardID)
values (?, ?);

-- Func: Get
SELECT *
FROM threads
WHERE threadID = ?;

-- Func: ByBoard
SELECT *
FROM threads
WHERE boardID = ?;

-- Func: ByPost
SELECT *
FROM threads