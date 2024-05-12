import psycopg2
from psycopg2 import sql
import dataclasses
from typing import Optional


@dataclasses.dataclass
class Client:
    user: str = "user"
    dbname: str = "kavuny"
    password: str = ""
    host: str = "localhost"
    port: int = 5432

    def __post_init__(self):
        self.conn = psycopg2.connect(
            user=self.user,
            dbname=self.dbname,
            password=self.password,
            host=self.host,
            port=self.port
        )

    def get_coffee_shop(self, id: Optional[int] = None):
        if id is None:
            query = sql.SQL("SELECT * FROM coffee_shops")
        else:
            query = sql.SQL("SELECT * FROM coffee_shops WHERE id = {}").format(
                sql.Placeholder()
            )
        cursor = self.conn.cursor()
        cursor.execute(query, [id])
        return cursor.fetchall()

    def insert_coffee_shop(self, shop: dict):
        cursor = self.conn.cursor()
        cursor.execute("SELECT nextval('coffee_shops_id_seq')")
        shop_id = cursor.fetchone()[0]

        query = sql.SQL("INSERT INTO coffee_shops VALUES ({})").format(
            sql.SQL(', ').join(sql.Placeholder() * (len(shop) + 1))
        )
        values = [v if not isinstance(v, list) else ",".join([str(i) for i in v])
                  for v in shop.values()]
        values.insert(0, shop_id)
        cursor.execute(query, values)
        self.conn.commit()
        return shop_id

    def update_coffee_shop(self, id: int, shop: dict):
        print(shop)
        query = sql.SQL("UPDATE coffee_shops SET {} WHERE id = %s").format(
            sql.SQL(', ').join(
                sql.SQL("{} = %s").format(sql.Identifier(k)) for k in shop.keys()
            )
        )
        cursor = self.conn.cursor()
        cursor.execute(query, list(shop.values()) + [id])
        self.conn.commit()
        return cursor.rowcount

    def get_shop_menu(self, id: int):
        query = sql.SQL("SELECT * FROM menus WHERE coffee_shop_id = {}").format((sql.Placeholder()))
        cursor = self.conn.cursor()
        cursor.execute(query, [id])
        return cursor.fetchall()

    def add_menu_item(self, id: int, item: dict):
        cursor = self.conn.cursor()
        cursor.execute("SELECT nextval('menus_id_seq')")
        item_id = cursor.fetchone()[0]
        query = sql.SQL("INSERT INTO menus VALUES ({})").format(
            sql.SQL(', ').join(sql.Placeholder() * (len(item) + 2))
        )
        values = [item_id, id]
        values.extend(item.values())
        print(values)
        cursor = self.conn.cursor()
        cursor.execute(query, values)
        self.conn.commit()
        return cursor.rowcount

    def delete_menu_item(self, id: int, item_id: int):
        query = sql.SQL("DELETE FROM menus WHERE coffee_shop_id = {} AND coffee_pack_id = {}")\
            .format(sql.Placeholder(), sql.Placeholder())
        cursor = self.conn.cursor()
        cursor.execute(query, (id, item_id))
        self.conn.commit()
        return cursor.rowcount
