from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from app.auth.dependencies import get_current_user
from app.auth.models import User
from app.music.dao import PlaylistDAO, SongDAO, RecentlyPlayedDAO
from app.music.schemas import PlaylistCreate, PlaylistRead, SongRead, RecentlyPlayedRead
from typing import List
from app.utils.dao.session_maker import SessionDep

router = APIRouter(prefix="/music", tags=["Music"])

@router.get("/playlists/", response_model=List[PlaylistRead])
async def get_playlists(
    session: AsyncSession = SessionDep,
    current_user: User = Depends(get_current_user)
):
    """Get all playlists for the current user"""
    return await PlaylistDAO.find_all(session=session, filters={"user_id": current_user.id})

@router.post("/playlists/", response_model=PlaylistRead)
async def create_playlist(
    playlist: PlaylistCreate,
    session: AsyncSession = SessionDep,
    current_user: User = Depends(get_current_user)
):
    """Create a new playlist"""
    playlist_dict = playlist.model_dump()
    playlist_dict["user_id"] = current_user.id
    return await PlaylistDAO.add(session=session, values=playlist_dict)

@router.get("/recently-played/", response_model=List[RecentlyPlayedRead])
async def get_recently_played(
    limit: int = 5,
    session: AsyncSession = SessionDep,
    current_user: User = Depends(get_current_user)
):
    """Get recently played tracks"""
    return await RecentlyPlayedDAO.find_all(
        session=session,
        filters={"user_id": current_user.id},
        limit=limit,
        order_by=[("-played_at", True)]
    )

@router.post("/recently-played/{song_id}")
async def add_recently_played(
    song_id: int,
    session: AsyncSession = SessionDep,
    current_user: User = Depends(get_current_user)
):
    """Add a song to recently played"""
    song = await SongDAO.find_one_or_none_by_id(song_id, session)
    if not song:
        raise HTTPException(status_code=404, detail="Song not found")
    
    await RecentlyPlayedDAO.add(session=session, values={
        "user_id": current_user.id,
        "song_id": song_id,
    })
    return {"message": "Song added to recently played"}

@router.get("/search/", response_model=List[SongRead])
async def search_songs(
    query: str,
    session: AsyncSession = SessionDep
):
    """Search for songs by title or artist name"""
    return await SongDAO.find_all(session=session, filters={"title__icontains": query})