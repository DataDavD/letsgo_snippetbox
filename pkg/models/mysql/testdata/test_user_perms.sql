CREATE USER 'test_web'@'localhost';

GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE ON test_snippetbox.* TO 'test_web'@'localhost';

ALTER USER 'test_web'@'localhost' IDENTIFIED BY 'pass';
