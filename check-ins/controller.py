import check_ins_service
from check_in_model import CheckIn
from fastapi import FastAPI, HTTPException

app = FastAPI()


@app.get("/check_ins")
def get_check_ins(
    coffee_shop_id: int | None = None,
    user_id: int | None = None,
    coffee_pack_id: int | None = None
):
    check_ins = check_ins_service.get_check_ins(coffee_shop_id, user_id, coffee_pack_id)
    if len(check_ins) == 0:
        raise HTTPException(status_code=404, detail="No check-ins found for the given parameters")
    return check_ins


@app.post("/check_ins")
def post_check_ins(check_in: CheckIn):
    # TO-DO add auth here
    user_id = -1
    check_in = check_in.model_dump()
    check_in["user_id"] = user_id
    return check_ins_service.post_check_ins(check_in)
