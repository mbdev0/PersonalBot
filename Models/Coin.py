from pydantic import BaseModel
from typing import Optional

class Coin(BaseModel):
    name: str
    symbol: str
    ipfs_url: str
    telegram_url: Optional[str] = None
    twitter_url: Optional[str] = None
    website_url: Optional[str] = None
    image_url: Optional[str] = None
