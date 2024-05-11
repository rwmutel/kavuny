CREATE TABLE IF NOT EXISTS menus (
    id SERIAL PRIMARY KEY,
    coffee_shop_id INT,
    coffee_pack_id INT,
    quantity INT,
    price FLOAT
);

COPY menus(id, coffee_shop_id, coffee_pack_id, quantity, price)

FROM '/opt/demo_data/menus.csv' DELIMITER E'\t' NULL AS 'null' CSV HEADER;
