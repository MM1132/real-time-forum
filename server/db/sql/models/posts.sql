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
       )                  AS likes,
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


-- Func: SearchPosts{p.postID,l2.date,likes}{ASC,DESC}
-- WITH COUNT
WITH opt AS
         (
             SELECT ? as myUserID,
                    ? as page,
                    ? as pageLen,

                    ? as threadID,
                    ? as threadTitle,
                    ? as tagName,
                    ? as author,
                    ? as authorID,
                    ? as boardID,
                    ? as boardName,

                    ? as likedBy,
                    ? as likedByID,
                    ? as dislikedBy,
                    ? as dislikedByID,

                    ? as content,

                    ? as after,
                    ? as before
         )

-- THIS SELECT
SELECT p.*,
       u.*,
       (
           SELECT IFNULL(SUM(value), 0)
           FROM likes
           WHERE postID = p.postID
       )                  AS likes,
       IFNULL(l.value, 0) AS myLikes,
       th.title,
       b.boardID,
       b.name
-- UNTIL FROM
FROM posts p
         LEFT JOIN
     users u on p.userID = u.userID
         LEFT JOIN
     threads th on p.threadID = th.threadID
         LEFT JOIN
     tags ta on th.threadID = ta.threadID
         LEFT JOIN
     boards b on th.boardID = b.boardID
         LEFT JOIN
     likes l on p.postID = l.postID AND l.userID = (select myUserID from opt)
         LEFT JOIN
     likes l2 on p.postID = l2.postID

WHERE ((select threadID from opt) IS NULL OR (select threadID from opt) = p.threadID)
  AND ((select author from opt) IS NULL OR (select author from opt) == u.name)
  AND ((select authorID from opt) IS NULL OR (select authorID from opt) == u.userID)
  AND ((select content from opt) IS NULL OR p.content LIKE '%' || (select content from opt) || '%')

  AND ((select boardID from opt) IS NULL OR (select boardID from opt) = th.boardID)
  AND ((select boardName from opt) IS NULL OR (select boardName from opt) = b.name)
  AND ((select likedBy from opt) IS NULL OR ((select u.userID from opt JOIN users u WHERE u.name = opt.likedBy LIMIT 1) = l2.userID AND l2.value = 1))
  AND ((select likedByID from opt) IS NULL OR ((select likedByID from opt) = l2.userID AND l2.value = 1))
  AND ((select dislikedBy from opt) IS NULL OR ((select u.userID from opt JOIN users u WHERE u.name = opt.dislikedBy LIMIT 1) = l2.userID AND l2.value = -1))
  AND ((select dislikedByID from opt) IS NULL OR ((select dislikedByID from opt) = l2.userID AND l2.value = -1))
  AND ((select threadTitle from opt) IS NULL OR th.title LIKE '%' || (select threadTitle from opt) || '%')
  AND ((select tagName from opt) IS NULL OR (select tagName from opt) = ta.name)

  AND ((select after from opt) IS NULL OR (select after from opt) < p.date)
  AND ((select before from opt) IS NULL OR (select before from opt) > p.date)

GROUP BY p.postID
ORDER BY
-- START
p.postID DESC
-- END
LIMIT (select pageLen from opt) OFFSET ((select page from opt) - 1) * (select pageLen from opt);
