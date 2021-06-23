USE snippetbox;

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
