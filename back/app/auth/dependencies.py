from fastapi import Request, HTTPException, status, Depends
from jose import jwt, JWTError
from datetime import datetime, timezone
from sqlalchemy.ext.asyncio import AsyncSession
from app.auth.utils import get_admin
from app.config import settings
from app.auth.dao import UsersDAO
from app.utils.dao.session_maker import SessionDep


def get_token(request: Request):
    token = request.cookies.get('users_access_token')
    if not token:
        raise Exception("Token not found in cookies")
    return token

async def get_current_user(token: str = Depends(get_token), session: AsyncSession = SessionDep):
    try:
        payload = jwt.decode(token, settings.SECRET_KEY, algorithms=settings.ALGORITHM)
    except JWTError:
        raise Exception("Invalid JWT token")

    expire: str = payload.get('exp')
    expire_time = datetime.fromtimestamp(int(expire), tz=timezone.utc)
    if not expire or expire_time < datetime.now(timezone.utc):
        raise Exception("Token has expired")

    user_id: str = payload.get('sub')
    if not user_id:
        raise Exception("User ID not found in token")

    user = await UsersDAO.find_one_or_none_by_id(data_id=int(user_id), session=session)
    if not user:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail='User not found')
    return user

if settings.DISABLE_AUTH == "true":
    get_current_user = get_admin
