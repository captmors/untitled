from datetime import datetime
from sqlalchemy import ForeignKey, String, Integer
from sqlalchemy.orm import Mapped, mapped_column, relationship
from app.utils.dao.database import Base
from app.auth.models import User

class Artist(Base):
    name: Mapped[str]
    bio: Mapped[str | None]
    image_url: Mapped[str | None]
    songs: Mapped[list["Song"]] = relationship(back_populates="artist")

class Song(Base):
    title: Mapped[str]
    duration: Mapped[int]  # Duration in seconds
    image_url: Mapped[str | None]
    artist_id: Mapped[int] = mapped_column(ForeignKey("artists.id"))
    artist: Mapped["Artist"] = relationship(back_populates="songs")
    playlists: Mapped[list["PlaylistSong"]] = relationship(back_populates="song")

class Playlist(Base):
    name: Mapped[str]
    description: Mapped[str | None]
    image_url: Mapped[str | None]
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"))
    user: Mapped["User"] = relationship()
    songs: Mapped[list["PlaylistSong"]] = relationship(back_populates="playlist")

class PlaylistSong(Base):
    playlist_id: Mapped[int] = mapped_column(ForeignKey("playlists.id"))
    song_id: Mapped[int] = mapped_column(ForeignKey("songs.id"))
    added_at: Mapped[datetime]
    playlist: Mapped["Playlist"] = relationship(back_populates="songs")
    song: Mapped["Song"] = relationship(back_populates="playlists")

class RecentlyPlayed(Base):
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"))
    song_id: Mapped[int] = mapped_column(ForeignKey("songs.id"))
    played_at: Mapped[datetime]
    user: Mapped["User"] = relationship()
    song: Mapped["Song"] = relationship()