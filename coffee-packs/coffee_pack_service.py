from persistence import Client
from typing import Optional, Tuple, List
from coffee_pack_model import CoffeePack

coffee_pack_item_keys = CoffeePack.model_fields.keys()
conn = Client("user", "kavuny", host="packs-db", port=5432)


def get_pack(id: int):
    return tuples_to_json(conn.get_coffee_packs([id]))


def get_packs(ids: Optional[list[int]] = None):
    return tuples_to_json(conn.get_coffee_packs(ids))


def tuples_to_json(
        tuple: List[Tuple],
        item_keys: list[str] = ["id"] + list(coffee_pack_item_keys)
        ):
    return [dict(zip(item_keys, item)) for item in tuple]


def create_pack(pack: CoffeePack) -> int:
    # TO-DO: add authentication before inserting the pack
    return conn.insert_coffee_pack(pack.model_dump())
