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

-- Func: ThreadCount
SELECT count(*)
FROM threads
WHERE boardID = ?;

-- Func: GetPageThreads{pl.postID,t.Title,CountPosts}{ASC,DESC}
WITH threadExtras AS
         (
             SELECT t.threadID,
                    count(*)                 AS countPosts,
                    count(distinct p.userID) AS countUsers,
                    max(p.postID)            AS latestID,
                    min(p.postID)            AS oldestID
             FROM threads t
                      JOIN
                  posts p ON t.threadID = p.threadID
             WHERE boardID = ?
             GROUP BY t.threadID
         )

SELECT t.*,
       countPosts,
       countUsers,

       latestID,
       ul.userID AS latestAuthorID,
       ul.name   AS latestAuthor,
       pl.date   AS latestDate,

       oldestID,
       uf.userID AS oldestAuthorID,
       uf.name   AS oldestAuthor,
       pf.date   AS oldestDate
FROM threads t
         JOIN
     threadExtras e ON t.threadID = e.threadID
         JOIN
     posts pl ON e.latestID = pl.postID
         JOIN
     users ul ON pl.userID = ul.userID
         JOIN
     posts pf ON e.oldestID = pf.postID
         JOIN
     users uf ON pf.userID = uf.userID
ORDER BY
--START
pl.postID DESC
--END
LIMIT ? OFFSET ?;


-- Func: GetPageThreadsByTag{pl.postID,t.Title,CountPosts}{ASC,DESC}
WITH threadExtras AS
         (
             SELECT t.threadID,
                    count(*)                 AS countPosts,
                    count(distinct p.userID) AS countUsers,
                    max(p.postID)            AS latestID,
                    min(p.postID)            AS oldestID
             FROM tags t
                      JOIN
                  posts p ON t.threadID = p.threadID
             WHERE t.name = ?
             GROUP BY t.threadID
         )

SELECT t.*,
       countPosts,
       countUsers,

       latestID,
       ul.userID AS latestAuthorID,
       ul.name   AS latestAuthor,
       pl.date   AS latestDate,

       oldestID,
       uf.userID AS oldestAuthorID,
       uf.name   AS oldestAuthor,
       pf.date   AS oldestDate
FROM threads t
         JOIN
     threadExtras e ON t.threadID = e.threadID
         JOIN
     posts pl ON e.latestID = pl.postID
         JOIN
     users ul ON pl.userID = ul.userID
         JOIN
     posts pf ON e.oldestID = pf.postID
         JOIN
     users uf ON pf.userID = uf.userID
ORDER BY
--START
pl.postID DESC
--END
LIMIT ? OFFSET ?;
