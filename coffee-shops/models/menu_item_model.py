from pydantic import BaseModel


class MenuItem(BaseModel):
    coffee_pack_id: int
    price: float
    quantity: int
