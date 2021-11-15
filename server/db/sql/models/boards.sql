-- Func: Insert
INSERT INTO boards(parentID, name, description)
values (?, ?, ?);

-- Func: Get
SELECT *
FROM boards
WHERE boardID = ?;

-- Func: GetChildren
SELECT *
FROM boards
WHERE parentID = ?
ORDER BY isGroup, "order";


-- Func: GetBreadcrumbs
WITH ancestors AS (
    SELECT *
    FROM boards
    WHERE boardID = ?

    UNION ALL

    SELECT b.*
    FROM boards b
             JOIN
         ancestors a ON b.boardID = a.parentID
)
SELECT *
FROM ancestors;

-- Func: GetExtras
WITH descendants AS
         (SELECT ? boardID

          UNION ALL

          SELECT b.boardID
          FROM boards b
                   JOIN
               descendants d
               ON d.boardID = b.parentID),

     aggregate AS
         (SELECT max(p.postID)              AS latestID,
                 count(distinct t.threadID) AS threadCount,
                 count(*)                   AS postCount
          FROM descendants d
                   JOIN
               threads t ON d.boardID = t.boardID
                   JOIN
               posts p ON t.threadID = p.threadID)


SELECT a.threadCount,
       a.postCount,

       p.postID   AS latestID,
       p.userID   AS latestAuthorID,
       u.name     AS latestAuthor,
       p.date     AS latestDate,

       t.threadID AS threadID,
       t.title    AS threadTitle
FROM aggregate a
         JOIN
     posts p ON a.latestID = p.postID
         JOIN
     threads t on t.threadID = p.threadID
         JOIN
     users u on u.userID = p.userID;
