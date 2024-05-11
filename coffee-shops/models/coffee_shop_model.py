from pydantic import BaseModel
from typing import List, Optional


class CoffeeShop(BaseModel):
    id: int
    name: str
    description: str
    image_path: str
    address_text: str
    address_latitude: float
    address_longitude: float
