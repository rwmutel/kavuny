# huge thanks to https://stackoverflow.com/questions/43598002/how-to-run-consul-on-docker-with-initial-key-value-pair-data

FROM consul:1.15.4

COPY . /opt/
RUN chmod 755 /opt/*

CMD /opt/start.sh

