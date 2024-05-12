import dataclasses
from typing import ClassVar, Optional

import psycopg2
from psycopg2 import sql


@dataclasses.dataclass
class Client:
    user: str = "user"
    dbname: str = "kavuny"
    password: str = ""
    host: str = "localhost"
    port: int = 5432
    get_packs_query: ClassVar[str] =\
        sql.SQL("SELECT * FROM coffee_packs WHERE id IN ({})")

    def __post_init__(self):
        self.conn = psycopg2.connect(
            user=self.user,
            dbname=self.dbname,
            password=self.password,
            host=self.host,
            port=self.port
        )

    def get_coffee_packs(self, ids: Optional[list[int]] = None):
        if ids:
            query = self.get_packs_query.format(
                sql.SQL(', ').join(sql.Placeholder() * len(ids))
            )
        else:
            query = sql.SQL("SELECT * FROM coffee_packs")
        cursor = self.conn.cursor()
        cursor.execute(query, ids)
        return cursor.fetchall()

    def insert_coffee_pack(self, pack: dict):
        cursor = self.conn.cursor()
        cursor.execute("SELECT nextval('coffee_packs_id_seq')")
        pack_id = cursor.fetchone()[0]

        query = sql.SQL("INSERT INTO coffee_packs VALUES ({})").format(
            sql.SQL(', ').join(sql.Placeholder() * (len(pack) + 1))
        )
        values = [v if not isinstance(v, list)
                  else ",".join([str(i) for i in v])
                  for v in pack.values()]
        values.insert(0, pack_id)
        cursor.execute(query, values)
        self.conn.commit()
        return pack_id
