-- create the databases
CREATE DATABASE IF NOT EXISTS blog;

-- create the users for each database
CREATE USER 'app'@'%' IDENTIFIED BY 'codev123';
GRANT CREATE, ALTER, INDEX, LOCK TABLES, REFERENCES, UPDATE, DELETE, DROP, SELECT, INSERT ON `blog`.* TO 'app'@'%';

FLUSH PRIVILEGES;