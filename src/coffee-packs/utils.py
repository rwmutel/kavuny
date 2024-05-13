import os
import random
import socket
import json
from datetime import datetime
from kafka.producer import KafkaProducer

import consul

c: consul.Consul = None
HTTP_PREFIX = "http://"


class UserType:
    USER = "user"
    SHOP = "shop"


def register_in_consul(name: str, port: int = 8080):
    global c
    consul_addr = os.getenv("CONSUL_ADDR").split(":")
    c = consul.Consul(host=consul_addr[0], port=int(consul_addr[1]))
    addr = socket.gethostbyname(socket.gethostname())
    service_id = name + "_" + addr
    c.agent.service.register(
        name=name,
        service_id=service_id,
        address=addr,
        port=port,
        check=consul.Check().http(
            f"{HTTP_PREFIX}{addr}:{port}/healthcheck", "30s"
            )
    )
    return service_id


def deregister_from_consul(service_id: str):
    c.agent.service.deregister(service_id)


def get_consul_kv(key: str):
    return c.kv.get(key)[1]["Value"].decode("utf-8")


def get_random_service_addr(name: str):
    service = random.choice(c.health.service(name)[1])["Service"]
    return HTTP_PREFIX + service["Address"] + ":" + str(service["Port"])


def log(data: dict):
    data["timestamp"] = datetime.now().timestamp()
    producer = KafkaProducer(bootstrap_servers=get_consul_kv("kafka_address"),
                             value_serializer=str.encode)
    producer.send(get_consul_kv("kafka_topic"),
                  json.dumps(data),
                  timestamp_ms=int(datetime.now().timestamp() * 1000))
