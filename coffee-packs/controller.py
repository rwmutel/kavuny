from fastapi import FastAPI, HTTPException, Query, Cookie
from typing import Annotated, Optional, List
import coffee_pack_service as service
from coffee_pack_model import CoffeePack
from contextlib import asynccontextmanager
import utils


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
    print(f"{session_id=}")
    try:
        pack_id = service.create_pack(pack, session_id)
        return f"Pack created with id {pack_id}"
    except Exception as e:
        raise e


@app.get("/healthcheck")
def healthcheck():
    return "OK"
