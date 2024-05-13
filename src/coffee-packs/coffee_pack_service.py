from typing import List, Optional, Tuple

import requests
import utils
from coffee_pack_model import CoffeePack
from fastapi import HTTPException
from persistence import Client

coffee_pack_item_keys = CoffeePack.model_fields.keys()
db_conf = utils.get_consul_kv("packs_db")
conn = Client(**db_conf)


def get_pack(id: int):
    return tuples_to_json(conn.get_coffee_packs([id]))


def get_packs(ids: Optional[list[int]] = None):
    return tuples_to_json(conn.get_coffee_packs(ids))


def tuples_to_json(
        tuple: List[Tuple],
        item_keys: list[str] = ["id"] + list(coffee_pack_item_keys)
        ):
    return [dict(zip(item_keys, item)) for item in tuple]


def create_pack(pack: CoffeePack, session_id: int) -> int:
    r = requests.get(utils.get_random_service_addr("auth-service") + "/id",
                     cookies={"session_id": session_id})
    if r.status_code == 401:
        raise HTTPException(status_code=401, detail=r.text)
    user_type, user_id = r.text.split(":")
    if user_type != utils.UserType.SHOP:
        raise HTTPException(status_code=401,
                            detail="Only shops can create packs")
    data = pack.model_dump().copy()
    data["author_id"] = user_id
    utils.log(data)
    return conn.insert_coffee_pack(pack.model_dump())
