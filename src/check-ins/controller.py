from contextlib import asynccontextmanager
from typing import Annotated

import check_ins_service
import utils
from check_in_model import CheckIn
from fastapi import Cookie, FastAPI, HTTPException


@asynccontextmanager
async def lifespan(app: FastAPI):
    service_id = utils.register_in_consul("check-ins")
    yield
    utils.deregister_from_consul(service_id)


app = FastAPI(lifespan=lifespan)


@app.get("/check-ins")
def get_check_ins(
    coffee_shop_id: int | None = None,
    user_id: int | None = None,
    coffee_pack_id: int | None = None
):
    check_ins = check_ins_service.get_check_ins(coffee_shop_id,
                                                user_id,
                                                coffee_pack_id)
    if len(check_ins) == 0:
        raise HTTPException(status_code=404,
                            detail="No check-ins for the given parameters")
    return check_ins


@app.post("/check-ins")
def post_check_ins(
    check_in: CheckIn,
    session_id: Annotated[str | None, Cookie()] = None
):
    return check_ins_service.post_check_ins(check_in.model_dump(), session_id)


@app.get("/healthcheck")
def healthcheck():
    return "OK"
