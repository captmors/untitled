from app.utils.dao.base import BaseDAO
from app.music.models import Artist, Song, Playlist, PlaylistSong, RecentlyPlayed

class ArtistDAO(BaseDAO):
    model = Artist

class SongDAO(BaseDAO):
    model = Song

class PlaylistDAO(BaseDAO):
    model = Playlist

class PlaylistSongDAO(BaseDAO):
    model = PlaylistSong

class RecentlyPlayedDAO(BaseDAO):
    model = RecentlyPlayed