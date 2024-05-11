CREATE TABLE IF NOT EXISTS coffee_packs (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    roastery VARCHAR,
    description TEXT,
    image_path VARCHAR,
    country VARCHAR,
    weight VARCHAR,
    flavour VARCHAR
);

COPY coffee_packs(id, name, roastery, description, image_path, country, weight, flavour)

FROM '/opt/demo_data/coffee_packs.csv' DELIMITER E'\t' NULL AS 'null' CSV HEADER;
