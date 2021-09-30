USE test_snippetbox;

CREATE TABLE snippets
(
    id      INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title   VARCHAR(100) NOT NULL,
    content TEXT         NOT NULL,
    created DATETIME     NOT NULL,
    expires DATETIME     NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets (created);

CREATE TABLE users
(
    id              INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    hashed_password CHAR(60)     NOT NULL,
    created         DATETIME     NOT NULL,
    active          BOOLEAN      NOT NULL DEFAULT TRUE
);

ALTER TABLE users
    ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created)
VALUES ('Alice Jones2',
        'alice2@example.com',
        '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
        '2018-12-23 17:25:22');
