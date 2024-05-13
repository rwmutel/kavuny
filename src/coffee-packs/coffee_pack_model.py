from typing import List, Optional

from pydantic import BaseModel


class CoffeePack(BaseModel):
    name: str
    roastery: str
    description: Optional[str] = None
    image_path: str
    country: str
    weight: List[int]
    flavour: List[str]
