from contextlib import asynccontextmanager
from typing import Annotated, List, Optional

import coffee_shop_service
import utils
from fastapi import Cookie, FastAPI, HTTPException
from models.menu_item_model import MenuItem


@asynccontextmanager
async def lifespan(app: FastAPI):
    service_id = utils.register_in_consul("coffee-shops")
    yield
    utils.deregister_from_consul(service_id)


app = FastAPI(lifespan=lifespan)


@app.get("/coffee-shops/{id}")
def get_coffee_shop(id: Optional[int] = None):
    shops = coffee_shop_service.get_coffee_shop(id)
    if not shops:
        raise HTTPException(status_code=404,
                            detail=f"Coffee shop with id {id} not found")
    return shops


@app.get("/coffee-shops/")
def get_all_shops():
    shops = coffee_shop_service.get_coffee_shop()
    if not shops:
        raise HTTPException(status_code=404, detail="No coffee shops found")
    return shops


@app.put("/coffee-shops/{id}")
def update_coffee_shop(
    id: int,
    coffee_shop: dict,
    session_id: Annotated[str | None, Cookie()] = None
):
    if coffee_shop_service.update_coffee_shop(id, coffee_shop, session_id):
        return "successfully updated"
    else:
        raise HTTPException(status_code=404,
                            detail=f"Coffee shop with id {id} not found")


@app.get("/coffee-shops/{id}/menu")
def get_shop_menu(id: int):
    menu = coffee_shop_service.get_shop_menu(id)
    if not menu:
        raise HTTPException(status_code=404,
                            detail="Menu for coffee shop "
                                   f"with id {id} not found")
    return menu


@app.post("/coffee-shops/{id}/menu")
def add_menu_item(
    id: int,
    item: MenuItem | List[MenuItem],
    session_id: Annotated[str | None, Cookie()] = None
):
    added_rows = 0
    if isinstance(item, list):
        for i in item:
            added_rows += coffee_shop_service.add_menu_item(id, i, session_id)
    else:
        added_rows += coffee_shop_service.add_menu_item(id, item, session_id)
    return f"{added_rows} item(s) added to menu"


@app.delete("/coffee-shops/{id}/menu")
def delete_menu_item(
    id: int,
    item_id: int,
    session_id: Annotated[str | None, Cookie()] = None
):
    num_deleted = coffee_shop_service.delete_menu_item(id, item_id, session_id)
    return f"successfully deleted {num_deleted} items"


@app.get("/healthcheck")
def healthcheck():
    return "OK"
