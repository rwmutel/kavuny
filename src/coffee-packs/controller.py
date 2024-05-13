from contextlib import asynccontextmanager
from typing import Annotated, List, Optional

import coffee_pack_service as service
import utils
from coffee_pack_model import CoffeePack
from fastapi import Cookie, FastAPI, Query


@asynccontextmanager
async def lifespan(app: FastAPI):
    service_id = utils.register_in_consul("coffee-packs")
    yield
    utils.deregiter_from_consul(service_id)


app = FastAPI(lifespan=lifespan)


@app.get("/packs/{id}")
def get_single_pack(id: int):
    return service.get_pack(id)


@app.get("/packs/")
def get_packs(ids: Annotated[Optional[List[int]], Query()] = None):
    return service.get_packs(ids)


@app.post("/packs/")
def create_pack(
    pack: CoffeePack,
    session_id: Annotated[str | None, Cookie()] = None
):
    pack_id = service.create_pack(pack, session_id)
    return f"Pack created with id {pack_id}"


@app.get("/healthcheck")
def healthcheck():
    return "OK"
