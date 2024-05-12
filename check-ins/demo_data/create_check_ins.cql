CREATE KEYSPACE IF NOT EXISTS kavuny WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 2};

USE kavuny;

CREATE TABLE shop_check_ins (
    user_id INT,
    coffee_shop_id INT,
    check_in_time TIMESTAMP,
    coffee_pack_id INT,
    rating INT,
    check_in_text TEXT,
    image_url TEXT,
    PRIMARY KEY ((coffee_shop_id), check_in_time, user_id)
);

CREATE TABLE pack_check_ins (
    user_id INT,
    coffee_shop_id INT,
    check_in_time TIMESTAMP,
    coffee_pack_id INT,
    rating INT,
    check_in_text TEXT,
    image_url TEXT,
    PRIMARY KEY ((coffee_pack_id), check_in_time, user_id)
);

COPY shop_check_ins (user_id, coffee_shop_id, check_in_time, coffee_pack_id, rating, check_in_text, image_url) FROM '/opt/app/check_ins.csv' WITH HEADER = TRUE;
COPY pack_check_ins (user_id, coffee_shop_id, check_in_time, coffee_pack_id, rating, check_in_text, image_url) FROM '/opt/app/check_ins.csv' WITH HEADER = TRUE;
