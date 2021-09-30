-- add some dummy records
USE snippetbox;
INSERT INTO snippets (title, content, created, expires)
VALUES ('An old silent pond',
        'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));

INSERT INTO snippets (title, content, created, expires)
VALUES ('Over the wintry forest',
        'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));

INSERT INTO snippets (title, content, created, expires)
VALUES ('First autumn morning',
        'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));

INSERT INTO snippets (title, content, created, expires)
VALUES ('DataDavD Awesome Adventures in Life',
        'DataDavD has had an awesome, super, crazy, cool life!!!',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY));
