from fastapi import FastAPI, HTTPException, Query
from typing import Annotated, Optional, List
import coffee_pack_service as service
from coffee_pack_model import CoffeePack

app = FastAPI()


@app.get("/packs/{id}")
def get_single_pack(id: int):
    return service.get_pack(id)


@app.get("/packs/")
def get_packs(ids: Annotated[Optional[List[int]], Query()] = None):
    return service.get_packs(ids)


@app.post("/packs/")
def create_pack(pack: CoffeePack):
    try:
        pack_id = service.create_pack(pack)
        return f"Pack created with id {pack_id}"
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))
