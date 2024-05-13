from contextlib import asynccontextmanager
import check_ins_service
from typing import Annotated
from check_in_model import CheckIn
from fastapi import FastAPI, HTTPException, Cookie
import utils


@asynccontextmanager
async def lifespan(app: FastAPI):
    service_id = utils.register_in_consul("check-ins")
    yield
    utils.deregister_from_consul(service_id)


app = FastAPI(lifespan=lifespan)


@app.get("/check_ins")
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


@app.post("/check_ins")
def post_check_ins(
    check_in: CheckIn,
    session_id: Annotated[str | None, Cookie()] = None
):
    return check_ins_service.post_check_ins(check_in.model_dump(), session_id)


@app.get("/healthcheck")
def healthcheck():
    return "OK"