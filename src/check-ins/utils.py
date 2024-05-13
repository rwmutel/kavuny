import json
import os
import random
import socket
from datetime import datetime

from kafka.producer import KafkaProducer

import consul


class UserType:
    USER = "user"
    SHOP = "shop"


consul_addr = os.getenv("CONSUL_ADDR").split(":")
c: consul.Consul = consul.Consul(host=consul_addr[0], port=int(consul_addr[1]))
HTTP_PREFIX = "http://"


def register_in_consul(name: str, port: int = 8080):
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
    byte_value = c.kv.get(key)[1]["Value"]
    try:
        value = json.loads(byte_value)
    except json.JSONDecodeError:
        value = byte_value.decode("utf-8")
    return value


def get_random_service_addr(name: str):
    service = random.choice(c.health.service(name)[1])["Service"]
    return HTTP_PREFIX + service["Address"] + ":" + str(service["Port"])


def log_checkin(check_in: dict):
    check_in["check_in_time"] = check_in["check_in_time"].timestamp()
    producer = KafkaProducer(bootstrap_servers=get_consul_kv("kafka_address"),
                             value_serializer=str.encode)
    producer.send(get_consul_kv("kafka_topic"),
                  json.dumps(check_in),
                  timestamp_ms=int(datetime.now().timestamp() * 1000))
