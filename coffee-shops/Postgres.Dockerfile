FROM postgres:latest

COPY ./demo_data/create_coffee_shops.sql /docker-entrypoint-initdb.d/create_coffee_shops.sql
COPY ./demo_data/create_menus.sql /docker-entrypoint-initdb.d/create_menus.sql

RUN echo "SELECT setval('coffee_shops_id_seq', (SELECT MAX(id) FROM coffee_shops));" > /docker-entrypoint-initdb.d/update_sequence_shops.sql
RUN echo "SELECT setval('menus_id_seq', (SELECT MAX(id) FROM menus));" > /docker-entrypoint-initdb.d/update_sequence_menus.sql

