from cassandra.cluster import Cluster, Session
from cassandra.query import dict_factory
from typing import Dict, Literal
from datetime import datetime


class Client:
    CLUSTER_ADDR = "check-ins-db-1"
    session: Session = None

    def __init__(self, addr: str | None = None):
        if Client.session is None:
            if addr is None:
                addr = self.CLUSTER_ADDR
            self.cluster = Cluster([addr])
            Client.session = self.cluster.connect(keyspace="kavuny")
            Client.session.row_factory = dict_factory

    def get_check_ins(self,
                      coffee_shop_id: int | None,
                      user_id: int | None,
                      coffee_pack_id: int | None) -> Dict:
        where_clause = []
        table = "shop_check_ins"
        if coffee_shop_id is not None:
            where_clause.append(f"coffee_shop_id = {coffee_shop_id}")
        if coffee_pack_id is not None:
            table = "pack_check_ins"
            where_clause.append(f"coffee_pack_id = {coffee_pack_id}")
        if user_id is not None:
            where_clause.append(f"user_id = {user_id}")
        query = f"SELECT * FROM {table}"
        if (coffee_shop_id is not None and coffee_pack_id is not None):
            where_clause[-1] += " ALLOW FILTERING"

        if where_clause:
            query += " WHERE " + " AND ".join(where_clause)
        return Client.session.execute(query).all()

    def post_check_in(self, check_in: Dict) -> Literal[True]:
        check_in["check_in_time"] = int(check_in["check_in_time"].timestamp() * 1000)
        values = check_in.values()
        print(values)
        query = "INSERT INTO pack_check_ins (coffee_shop_id, check_in_time, coffee_pack_id, rating, check_in_text, user_id) " \
                "VALUES (%s, %s, %s, %s, %s, %s)"
        Client.session.execute(query, values)
        query = "INSERT INTO shop_check_ins (coffee_shop_id, check_in_time, coffee_pack_id, rating, check_in_text, user_id) " \
                "VALUES (%s, %s, %s, %s, %s, %s)"
        Client.session.execute(query, values)
        return True
