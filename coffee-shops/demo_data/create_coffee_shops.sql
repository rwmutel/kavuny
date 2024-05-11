CREATE TABLE IF NOT EXISTS coffee_shops (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    description TEXT,
    image_path VARCHAR,
    address_text VARCHAR,
    address_latitude FLOAT,
    address_longitude FLOAT
);

COPY coffee_shops(id, name, description, image_path, address_text, address_latitude, address_longitude)

FROM '/opt/demo_data/coffee_shops.csv' DELIMITER E'\t' NULL AS 'null' CSV HEADER;
