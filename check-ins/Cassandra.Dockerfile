FROM cassandra:latest

COPY cassandra_entry.sh /opt/app/cassandra_entry.sh
COPY demo_data/create_check_ins.cql /opt/app/create_check_ins.cql
COPY demo_data/check_ins.csv /opt/app/check_ins.csv

ENTRYPOINT ["bash", "/opt/app/cassandra_entry.sh"]
CMD ["cassandra", "-f"]
