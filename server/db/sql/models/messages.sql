-- Func: New
INSERT INTO messages(fromID, toID, body, date)
values (?, ?, ?, ?);

-- Func: GetByFrom
SELECT *
FROM messages
WHERE fromID = ?;

-- Func: GetByTo
SELECT *
FROM messages
WHERE toID = ?;

-- Func: GetChat
SELECT *
FROM messages
WHERE (fromID = ? OR fromID = ?) AND (toID = ? OR toID = ?)
ORDER BY MessageID DESC LIMIT ?, 10;