CREATE TABLE IF NOT EXISTS `users`
(
    `userID`      INTEGER PRIMARY KEY AUTOINCREMENT,
    `nickname`        TEXT UNIQUE NOT NULL,
    `firstname`   TEXT NOT NULL DEFAULT 'Robert',
    `lastname`   TEXT NOT NULL DEFAULT 'Reimann',
    `age`   DATE NOT NULL DEFAULT '0001-01-01',
    `gender`   TEXT NOT NULL DEFAULT 'Female',
    `email`       TEXT UNIQUE NOT NULL,
    `password`    TEXT        NOT NULL,
    `image`       TEXT        NOT NULL DEFAULT '0-0.png',
    `description` TEXT        NOT NULL DEFAULT 'Welcome to my profile!',
    `created`     DATE        NOT NULL
);

CREATE TABLE IF NOT EXISTS `sessions`
(
    `token`   TEXT PRIMARY KEY,
    `userID`  INTEGER NOT NULL,
    `created` DATE    NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (userID)
        ON DELETE CASCADE,
    -- Make sessions one-per-user, just for audit requirement
    CONSTRAINT session_unique UNIQUE (userID)
        ON CONFLICT REPLACE
);

CREATE TABLE IF NOT EXISTS `threads`
(
    `threadID` INTEGER PRIMARY KEY AUTOINCREMENT,
    `title`    TEXT    NOT NULL,
    `boardID`  INTEGER NOT NULL,
    FOREIGN KEY (boardID) REFERENCES boards (boardID)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS threads_boardID
    ON threads (boardID, threadID);

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

CREATE INDEX IF NOT EXISTS posts_threadID
    ON posts (threadID, postID);



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

CREATE TABLE IF NOT EXISTS `likes`
(
    `postID` INTEGER NOT NULL,
    `userID` INTEGER NOT NULL,
    `value`  INTEGER NOT NULL,
    `date`   DATE    NOT NULL,
    FOREIGN KEY (postID) REFERENCES posts (postID),
    FOREIGN KEY (userID) REFERENCES users (userID),
    CONSTRAINT uniq UNIQUE (postID, userID)
        ON CONFLICT REPLACE
);

CREATE INDEX IF NOT EXISTS likes_userID
    ON likes (userID, value);

CREATE TABLE IF NOT EXISTS `tags`
(
    `name`     TEXT,
    `threadID` INTEGER NOT NULL,
    FOREIGN KEY (threadID) REFERENCES threads (threadID)
        ON DELETE CASCADE,
    CONSTRAINT tag_unique UNIQUE (threadID, name)
        ON CONFLICT IGNORE
);

CREATE INDEX IF NOT EXISTS tag_index ON tags (name, threadID);

CREATE TABLE IF NOT EXISTS `pings`
(
    `pingID`  INTEGER PRIMARY KEY AUTOINCREMENT,
    `userID`  INTEGER NOT NULL,
    `content` TEXT    NOT NULL,
    `link`    TEXT    NOT NULL,
    `date`    DATE    NOT NULL,
    FOREIGN KEY (userID) REFERENCES users (userID)
);


CREATE TABLE IF NOT EXISTS `messages`
(
    `messageID` INTEGER PRIMARY KEY AUTOINCREMENT,
    `fromID`    INTEGER    NOT NULL,
    `toID`  INTEGER   NOT NULL,
    `body`  TEXT    NOT NULL,
    `date`  DATE    NOT NULL

);