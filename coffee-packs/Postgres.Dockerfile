FROM postgres:latest

COPY ./demo_data/create_coffee_packs.sql /docker-entrypoint-initdb.d/create_coffee_packs.sql
RUN echo "SELECT setval('coffee_packs_id_seq', (SELECT MAX(id) FROM coffee_packs));" > /docker-entrypoint-initdb.d/update_sequence.sql
