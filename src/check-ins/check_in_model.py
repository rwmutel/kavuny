from pydantic import BaseModel
from datetime import datetime


class CheckIn(BaseModel):
    coffee_shop_id: int | None
    check_in_time: datetime
    coffee_pack_id: int | None
    rating: int
    check_in_text: str | None
