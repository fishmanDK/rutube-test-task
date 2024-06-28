CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    hash_password VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS birth_subs
(
    root_user INT NOT NULL,
    subscriber INT NOT NULL,

    FOREIGN KEY (root_user) REFERENCES users(id),
    FOREIGN KEY (subscriber) REFERENCES users(id),
    UNIQUE (root_user, subscriber)
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    userID INT NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,

    FOREIGN KEY (userID) REFERENCES users(id)
);