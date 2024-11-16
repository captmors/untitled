from datetime import datetime
from pydantic import BaseModel, ConfigDict

class ArtistBase(BaseModel):
    name: str
    bio: str | None = None
    image_url: str | None = None
    model_config = ConfigDict(from_attributes=True)

class SongBase(BaseModel):
    title: str
    duration: int
    image_url: str | None = None
    artist_id: int
    model_config = ConfigDict(from_attributes=True)

class PlaylistBase(BaseModel):
    name: str
    description: str | None = None
    image_url: str | None = None
    model_config = ConfigDict(from_attributes=True)

class PlaylistCreate(PlaylistBase):
    pass

class PlaylistRead(PlaylistBase):
    id: int
    user_id: int
    created_at: datetime

class SongRead(SongBase):
    id: int
    artist: ArtistBase

class RecentlyPlayedRead(BaseModel):
    song: SongRead
    played_at: datetime
    model_config = ConfigDict(from_attributes=True)