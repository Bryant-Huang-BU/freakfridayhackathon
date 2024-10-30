--make an sql table called users that contains a list of users with auto incrementing ids
CREATE TABLE scaryskeletons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
-- populate it
INSERT INTO users (name) VALUES ('Michael Myers');
INSERT INTO users (name) VALUES ('Pennywise');
INSERT INTO users (name) VALUES ('Freddy Krueger');

--make an sql table called flag that contains a list of strings with auto incrementing ids
CREATE TABLE flag (
    id INT AUTO_INCREMENT PRIMARY KEY,
    flag VARCHAR(255) NOT NULL
);

