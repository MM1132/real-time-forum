-- Func: Insert
INSERT INTO threads(title, boardID)
values (?, ?);

-- Func: Get
SELECT *
FROM threads
WHERE threadID = ?;

-- Func: GetOP
SELECT *
FROM posts p
         JOIN users u ON p.userID = u.userID
WHERE threadID = ?
ORDER BY postID
LIMIT 1;

-- Func: ByBoard
SELECT *
FROM threads
WHERE boardID = ?;

-- Func: ThreadCount
SELECT count(*)
FROM threads
WHERE boardID = ?;

-- Func: SearchThreads{pl.postID,t.Title,CountPosts}{ASC,DESC}
-- WITH COUNT
WITH opt AS
         (
             SELECT ? as page,
                    ? as pageLen,

                    ? as threadTitle,
                    ? as tagName,
                    ? as author,
                    ? as boardID,
                    ? as boardName,

                    ? as latestAfter,
                    ? as latestBefore,
                    ? as oldestAfter,
                    ? as oldestBefore
         ),

     filterThreads AS
         (
             SELECT th.*
             FROM threads th
                      LEFT JOIN
                  boards b ON th.boardID = b.boardID
                      LEFT JOIN
                  tags ta ON th.threadID = ta.threadID
                      LEFT JOIN
                  posts p ON th.threadID = p.threadID

             WHERE ((select boardID from opt) IS NULL OR (select boardID from opt) = th.boardID)
               AND ((select boardName from opt) IS NULL OR (select boardName from opt) = b.name)
               AND ((select threadTitle from opt) IS NULL OR th.title LIKE '%' || (select threadTitle from opt) || '%')
               AND ((select tagName from opt) IS NULL OR (select tagName from opt) = ta.name)

             GROUP BY th.threadID
         ),

     threadExtras AS
         (
             SELECT t.threadID,
                    count(*)                 AS countPosts,
                    count(distinct p.userID) AS countUsers,
                    max(p.postID)            AS latestID,
                    min(p.postID)            AS oldestID
             FROM filterThreads t
                      JOIN
                  posts p ON t.threadID = p.threadID
             GROUP BY t.threadID
         )

-- THIS SELECT
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
-- UNTIL FROM
FROM filterThreads t
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
WHERE ((select latestAfter from opt) IS NULL OR (select latestAfter from opt) < pl.date)
  AND ((select latestBefore from opt) IS NULL OR (select latestBefore from opt) > pl.date)
  AND ((select oldestAfter from opt) IS NULL OR (select oldestAfter from opt) < pf.date)
  AND ((select oldestBefore from opt) IS NULL OR (select oldestBefore from opt) > pf.date)
  AND ((select author from opt) IS NULL OR (select author from opt) = uf.name)

ORDER BY
--START
pl.postID DESC
--END
LIMIT (select pageLen from opt) OFFSET ((select page from opt) - 1) * (select pageLen from opt);
