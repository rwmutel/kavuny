from typing import List, Optional

from pydantic import BaseModel


class CoffeeShop(BaseModel):
    id: int
    name: str
    description: str
    image_path: str
    address_text: str
    address_latitude: float
    address_longitude: float
