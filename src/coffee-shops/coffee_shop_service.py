from typing import List, Tuple

import pandas as pd
import requests
import utils
from fastapi import HTTPException
from models.coffee_shop_model import CoffeeShop
from models.menu_item_model import MenuItem
from persistence import Client

db_conf = utils.get_consul_kv("shops_db")
conn = Client(**db_conf)


def get_coffee_shop(id: int = None):
    return tuples_to_json(conn.get_coffee_shop(id))


def update_coffee_shop(id: int, coffee_shop: dict, session_id: int):
    uid = get_id_by_session(session_id)
    user_type, user_id = uid.split(":")
    if not (user_type == "shop" and int(user_id) == id):
        raise HTTPException(status_code=403,
                            detail="You are not allowed to update this shop")
    utils.log(coffee_shop.copy())
    return conn.update_coffee_shop(id, coffee_shop)


def get_shop_menu(id: int):
    shop_items = pd.DataFrame(conn.get_shop_menu(id),
                              columns=[
                                  "id",
                                  "coffee_shop_id",
                                  "pack_id",
                                  "quantity",
                                  "price"
                                  ])
    item_details = pd.DataFrame(
        requests.get(utils.get_random_service_addr("coffee-packs") + "/packs/",
                     params={"ids": list(shop_items["pack_id"])}).json()
    )
    return shop_items\
        .merge(item_details, left_on="pack_id", right_on="id")\
        .drop(columns=["id_x", "id_y"])\
        .to_dict(orient="records")


def add_menu_item(id: int, item: MenuItem, session_id: int):
    uid = get_id_by_session(session_id)
    user_type, user_id = uid.split(":")
    if not (user_type == "shop" and int(user_id) == id):
        raise HTTPException(status_code=403,
                            detail="You are not allowed "
                                   "to add items to this shop")
    utils.log(item.model_dump().copy())
    return conn.add_menu_item(id, item.model_dump())


def delete_menu_item(id: int, item_id: int, session_id: int):
    uid = get_id_by_session(session_id)
    user_type, user_id = uid.split(":")
    if not (user_type == "shop" and int(user_id) == id):
        raise HTTPException(status_code=403,
                            detail="You are not allowed to delete this item")
    data = {"action": "deleted item", "shop_id": id, "item_id": item_id}
    utils.log(data)
    return conn.delete_menu_item(id, item_id)


def tuples_to_json(
    tuple: List[Tuple],
    item_keys: list[str] = CoffeeShop.model_fields.keys()
):
    return [dict(zip(item_keys, item)) for item in tuple]


def get_id_by_session(session_id: str):
    r = requests.get(utils.get_random_service_addr("auth-service") + "/id",
                     cookies={"session_id": session_id})
    if r.status_code == 401:
        raise HTTPException(status_code=401, detail=r.text)
    return r.text
