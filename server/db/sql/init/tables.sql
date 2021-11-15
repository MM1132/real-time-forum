CREATE TABLE IF NOT EXISTS `users`
(
    `userID`   INTEGER PRIMARY KEY AUTOINCREMENT,
    `name`     TEXT UNIQUE NOT NULL,
    `email`    TEXT UNIQUE NOT NULL,
    `password` TEXT        NOT NULL,
    `created`  DATE        NOT NULL
);

CREATE TABLE IF NOT EXISTS `sessions`
(
    `token`   TEXT PRIMARY KEY,
    `userID`  INTEGER NOT NULL,
    `created` DATE    NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (userID)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `threads`
(
    `threadID` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title`    TEXT    NOT NULL,
    `boardID`  INTEGER NOT NULL,
    FOREIGN KEY (boardID) REFERENCES boards (boardID)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `posts`
(
    `postID`   INTEGER PRIMARY KEY AUTOINCREMENT,
    `threadID` INTEGER,
    `userID`   INTEGER,
    `content`  TEXT NOT NULL,
    `date`     DATE NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (userID),
    FOREIGN KEY (threadID) REFERENCES threads (threadID)
        ON DELETE CASCADE
);



CREATE TABLE IF NOT EXISTS `boards`
(
    `boardID`     INTEGER PRIMARY KEY AUTOINCREMENT,
    `parentID`    INTEGER,
    `name`        TEXT    NOT NULL,
    `description` TEXT,
    `isGroup`     INTEGER NOT NULL DEFAULT 0,
    `order`       INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (parentID) REFERENCES boards (boardID)
        ON DELETE RESTRICT
);
