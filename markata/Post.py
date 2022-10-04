from pydantic import BaseModel


class Post(BaseModel):
    title: str = None
    date: date = None
