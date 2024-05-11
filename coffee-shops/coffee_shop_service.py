from persistence import Client
from models.coffee_shop_model import CoffeeShop
from models.menu_item_model import MenuItem
from typing import List, Tuple
import requests
import pandas as pd

conn = Client("user", "kavuny", host="shops-db", port=5432)


def get_coffee_shop(id: int = None):
    return tuples_to_json(conn.get_coffee_shop(id))


def update_coffee_shop(id: int, coffee_shop: dict):
    # TO-DO: auth here
    return conn.update_coffee_shop(id, coffee_shop)


def get_shop_menu(id: int):
    shop_items = pd.DataFrame(conn.get_shop_menu(id),
                              columns=["id", "coffee_shop_id", "pack_id", "quantity", "price"])
    item_details = pd.DataFrame(
        requests.get("http://coffee-packs:8080/packs/",
                     params={"ids": list(shop_items["pack_id"])}).json()
    )
    return shop_items\
        .merge(item_details, left_on="pack_id", right_on="id")\
        .drop(columns=["id_x", "id_y"])\
        .to_dict(orient="records")


def add_menu_item(id: int, item: MenuItem):
    # TO-DO: auth here
    return conn.add_menu_item(id, item.model_dump())


def delete_menu_item(id: int, item_id: int):
    # TO-DO: auth here
    return conn.delete_menu_item(id, item_id)


def tuples_to_json(
        tuple: List[Tuple],
        item_keys: list[str] = CoffeeShop.model_fields.keys()
        ):
    return [dict(zip(item_keys, item)) for item in tuple]
