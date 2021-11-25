-- Func: Insert
INSERT INTO tags
VALUES (?, ?);

-- Func: Delete
DELETE
FROM tags
WHERE threadID = ?
  AND name = ?;

-- Func: GetThreads
SELECT th.*
FROM tags ta
         JOIN
     threads th ON th.threadID = ta.threadID
WHERE ta.name = ?
ORDER BY th.threadID DESC;

-- Func: GetByThread
SELECT tags.name
FROM tags
WHERE threadID = ?
ORDER BY ROWID;

-- Func: ThreadCount
SELECT count(*)
FROM tags
WHERE name = ?;

-- Func: GetPopular
SELECT tags.name
FROM threads
         JOIN
     tags on threads.threadID = tags.threadID
WHERE threads.boardID = ?
GROUP BY tags.name
ORDER BY count(*) DESC
LIMIT 100;
