-- Func: Insert
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

-- Func: GetByID
SELECT *
FROM messages
WHERE messageID = ?;

-- Func: GetRecentIDs
SELECT CASE WHEN fromID != $1 THEN fromID ELSE toID END recentID
FROM messages
WHERE fromID = $1 OR toID = $1
GROUP BY recentID
ORDER BY max(messageID) DESC;
