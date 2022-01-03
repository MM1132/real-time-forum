-- Func: Send
INSERT INTO pings(userID, content, link, date)
values (?, ?, ?, ?);

-- Func: Delete
DELETE FROM pings
WHERE pingID = ?;

-- Func: GetByUser
SELECT *
FROM pings
WHERE userID = ?;
