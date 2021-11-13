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
);

CREATE TABLE IF NOT EXISTS `threads`
(
    `threadID`   INTEGER PRIMARY KEY AUTOINCREMENT,
    `title`      TEXT    NOT NULL,
    `categoryID` INTEGER NOT NULL,
    FOREIGN KEY (categoryID) REFERENCES categories (categoryID)
);

CREATE TABLE IF NOT EXISTS `posts`
(
    `postID`   INTEGER PRIMARY KEY AUTOINCREMENT,
    `threadID` INTEGER,
    `userID`   INTEGER NOT NULL,
    `content`  TEXT    NOT NULL,
    `date`     DATE    NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (userID),
    FOREIGN KEY (threadID) REFERENCES threads (threadID)
);

CREATE TABLE IF NOT EXISTS `categories`
(
    `categoryID`  INTEGER PRIMARY KEY AUTOINCREMENT,
    `parentID`    INTEGER,
    `name`        TEXT NOT NULL,
    `description` TEXT,
    FOREIGN KEY (parentID) REFERENCES categories (categoryID)
);
