CREATE TABLE IF NOT EXISTS users(
    uid VARCHAR(36) PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS notes(
    nid VARCHAR(36) PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(uid),
    UNIQUE(title)
);
