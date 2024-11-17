from fastapi import Depends, HTTPException
from loguru import logger
from passlib.context import CryptContext
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from app.auth.models import Role, User
from app.utils.dao.session_maker import SessionDep, session_manager

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def get_password_hash(password: str) -> str:
    return pwd_context.hash(password)


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


async def get_admin(session: AsyncSession = SessionDep) -> User:
    result = await session.execute(select(User).join(Role).filter(Role.name == "admin"))
    admin_user = result.scalars().first()
    
    if not admin_user:
        raise HTTPException(status_code=404, detail="Admin user not found.")
    
    return admin_user

async def create_admin_role_and_user(session: AsyncSession):
    # Создаем роль администратора, если она не существует
    result = await session.execute(select(Role).filter_by(name="admin"))
    admin_role = result.scalars().first()
    
    if not admin_role:
        admin_role = Role(name="admin")
        session.add(admin_role)
        await session.commit()
        logger.info("Created admin role.")
    
    # Создаем пользователя с ролью администратора, если он не существует
    result = await session.execute(select(User).filter_by(email="admin@example.com"))
    existing_user = result.scalars().first()

    if not existing_user:
        admin_user = User(
            email="admin@example.com",
            first_name="Admin",
            last_name="Admin",
            phone_number="2281337",
            password=get_password_hash("admin"),
            role=admin_role
        )
        session.add(admin_user)
        await session.commit()
        logger.info("Created mock admin user with admin role.")
        return admin_user
    else:
        logger.info("Admin user already exists.")
        return existing_user
    