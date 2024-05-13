CREATE DATABASE auth;
CREATE TYPE USER_TYPES AS ENUM ('user', 'shop');
CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    login     VARCHAR NOT NULL,
    password  VARCHAR NOT NULL,
    salt      VARCHAR NOT NULL,
    user_type USER_TYPES NOT NULL
);
COPY users (id, login, password, salt, user_type)
    FROM '/opt/demo_data/users.csv' DELIMITER ',' NULL AS 'null' CSV HEADER;
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
