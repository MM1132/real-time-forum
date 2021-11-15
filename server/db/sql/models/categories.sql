-- Func: Insert
INSERT INTO categories(parentID, name, description)
values (?, ?, ?);

-- Func: Get
SELECT *
FROM categories
WHERE categoryID = ?;

-- Func: GetChildren
SELECT *
FROM categories
WHERE parentID = ?;

-- Func: GetBreadcrumbs
WITH ancestors AS (
    SELECT *
    FROM categories
    WHERE categoryID = ?

    UNION ALL

    SELECT c.*
    FROM categories c
             JOIN
         ancestors a ON c.categoryID = a.parentID
)
SELECT *
FROM ancestors;

-- Func: GetExtras
WITH descendants AS
         (SELECT ? categoryID

          UNION ALL

          SELECT c.categoryID
          FROM categories c
                   JOIN
               descendants d
               ON d.categoryID = c.parentID),

     aggregate AS
         (SELECT max(p.postID)              AS latestID,
                 count(distinct t.threadID) AS threadCount,
                 count(*)                   AS postCount
          FROM descendants d
                   JOIN
               threads t ON d.categoryID = t.categoryID
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
