from fastapi import HTTPException
from persistence import Client
import requests
import utils

conn = Client()


def get_check_ins(coffee_shop_id: int | None,
                  user_id: int | None,
                  coffee_pack_id: int | None):
    return conn.get_check_ins(coffee_shop_id, user_id, coffee_pack_id)


def post_check_ins(check_in: dict, session_id: str):
    r = requests.get(utils.get_random_service_addr("auth-service") + "/id",
                     cookies={"session_id": session_id})
    if r.status_code == 401:
        raise HTTPException(status_code=401, detail=r.text)
    user_type, user_id = r.text.split(":")
    if user_type != utils.UserType.USER:
        raise HTTPException(status_code=403,
                            detail="Only users can leave check-ins")
    check_in["user_id"] = int(user_id)
    utils.log_checkin(check_in.copy())
    return conn.post_check_in(check_in)
