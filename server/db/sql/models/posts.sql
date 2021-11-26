-- Func: Insert
INSERT INTO posts(content, userID, threadID, date)
values (?, ?, ?, ?);

-- Func: Get
SELECT *
FROM posts p
         JOIN users u ON p.userID = u.userID
WHERE postID = ?;

-- Func: GetByThreadID
SELECT p.*,
       u.*,
       (SELECT IFNULL(SUM(value), 0)
        FROM likes
        WHERE postID = p.postID
       ) AS likes,
       IFNULL(l.value, 0) AS myLikes


FROM posts p
         JOIN users u ON p.userID = u.userID
         LEFT JOIN likes l on p.postID = l.postID AND l.userID = ?
WHERE threadID = ?
ORDER BY postID;

-- Func: GetByUserID
SELECT *
FROM posts
WHERE userID = ?
ORDER BY postID DESC;
