-- Func: Insert
INSERT INTO posts(content, userID, threadID, date) values(?,?,?,?);

-- Func: Get
SELECT * FROM posts p JOIN users u ON p.userID = u.userID WHERE postID=?;

-- Func: GetByThreadID
SELECT * FROM posts p JOIN users u ON p.userID = u.userID WHERE threadID=? ORDER BY date;

-- Func: GetByUserID
SELECT * FROM posts WHERE userID=? ORDER BY date DESC;
