-- Func: Insert
INSERT INTO likes(postID, userID, value, date)
values (?, ?, ?, ?);

-- Func: GetPostTotal
SELECT IFNULL(SUM(value), 0)
FROM likes
WHERE postID = ?;

-- Func: GetAllUser
SELECT *
FROM likes
WHERE userID = ?
ORDER BY date DESC;

-- Func: Delete
DELETE
FROM likes
WHERE postID = ?
  AND userID = ?;

-- Func: Get
SELECT *
FROM likes
WHERE postID = ?
  AND userID = ?;

-- Func: Update
UPDATE likes
SET value = ?,
    date  = ?
WHERE userID = ?
  AND postID = ?;

-- Func: GetLikedPosts
SELECT l.value, p.*, author.*
FROM users u
         JOIN likes l on u.userID = l.userID
         JOIN posts p on p.postID = l.postID
         JOIN users author on author.userID = p.userID
WHERE u.userID = ?;
