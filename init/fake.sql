DROP TABLE IF EXISTS user;
CREATE TABLE user (
   username TEXT PRIMARY KEY NOT NULL,
   password TEXT NOT NULL,
   created TIMESTAMP
);


DROP TABLE IF EXISTS account;
CREATE TABLE account (
   username TEXT PRIMARY KEY NOT NULL,
   password TEXT NOT NULL,
   created TIMESTAMP,
   updated TIMESTAMP,
   website TEXT
);


INSERT INTO user(username, password) VALUES ('admin', 'admin');