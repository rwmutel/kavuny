CREATE DATABASE auth;
CREATE TYPE USER_TYPES AS ENUM ('user', 'shop');
CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    salt     VARCHAR NOT NULL,
    userType USER_TYPES
);