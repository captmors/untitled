from fastapi import APIRouter, File, UploadFile, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from app.utils.s3.utils import upload_file_to_minio
from app.auth.dependencies import get_current_user, SessionDep
from app.auth.models import User

router = APIRouter(prefix="/user", tags=["User"])

@router.post("/upload-avatar/")
async def upload_avatar(
    file: UploadFile = File(...),
    session: AsyncSession = SessionDep,
    current_user: User = Depends(get_current_user)
):
    """
    Загрузка аватара пользователя.
    """
    file_data = await file.read()
    avatar_url = upload_file_to_minio("user-avatars", file_data, file.filename)
    
    # Обновляем URL аватара в БД
    current_user.avatar_url = avatar_url
    session.add(current_user)
    await session.commit()
    
    return {"ok": True, "avatar_url": avatar_url}
