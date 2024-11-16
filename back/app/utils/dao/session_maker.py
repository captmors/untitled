from contextlib import asynccontextmanager
from typing import Callable, Optional, AsyncGenerator
from fastapi import Depends
from loguru import logger
from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession
from sqlalchemy import text
from functools import wraps
from .database import async_session_maker

class DatabaseSessionManager:
    """
    A class to manage asynchronous database sessions, including support for transactions and FastAPI dependencies.
    """

    def __init__(self, session_maker: async_sessionmaker[AsyncSession]):
        self.session_maker = session_maker

    @asynccontextmanager
    async def create_session(self) -> AsyncGenerator[AsyncSession, None]:
        """
        Creates and provides a new database session. Ensures the session is closed after use.
        """
        async with self.session_maker() as session:
            try:
                yield session
            except Exception as e:
                logger.error(f"Database session creation error: {e}")
                raise
            finally:
                await session.close()

    @asynccontextmanager
    async def transaction(self, session: AsyncSession) -> AsyncGenerator[None, None]:
        """
        Manages a transaction: commits on success, rolls back on error.
        """
        try:
            yield
            await session.commit()
        except Exception as e:
            await session.rollback()
            logger.exception(f"Transaction error: {e}")
            raise

    async def get_session(self) -> AsyncGenerator[AsyncSession, None]:
        """
        FastAPI dependency returning a session without transaction management.
        """
        async with self.create_session() as session:
            yield session

    async def get_transaction_session(self) -> AsyncGenerator[AsyncSession, None]:
        """
        FastAPI dependency returning a session with transaction management.
        """
        async with self.create_session() as session:
            async with self.transaction(session):
                yield session

    def connection(self, isolation_level: Optional[str] = None, commit: bool = True):
        """
        Decorator to manage a session with optional isolation level and commit configuration.

        Parameters:
        - `isolation_level`: transaction isolation level (e.g., "SERIALIZABLE").
        - `commit`: if `True`, commits after method execution.
        """

        def decorator(method):
            @wraps(method)
            async def wrapper(*args, **kwargs):
                async with self.session_maker() as session:
                    try:
                        if isolation_level:
                            await session.execute(text(f"SET TRANSACTION ISOLATION LEVEL {isolation_level}"))

                        result = await method(*args, session=session, **kwargs)

                        if commit:
                            await session.commit()

                        return result
                    except Exception as e:
                        await session.rollback()
                        logger.error(f"Transaction execution error: {e}")
                        raise
                    finally:
                        await session.close()

            return wrapper

        return decorator

    @property
    def session_dependency(self) -> Callable:
        """Returns a FastAPI dependency providing access to a session without transaction management."""
        return Depends(self.get_session)

    @property
    def transaction_session_dependency(self) -> Callable:
        """Returns a FastAPI dependency with transaction support."""
        return Depends(self.get_transaction_session)


# Initialize the session manager with the database session maker
session_manager = DatabaseSessionManager(async_session_maker)

# FastAPI dependencies for using sessions
SessionDep = session_manager.session_dependency
TransactionSessionDep = session_manager.transaction_session_dependency

# Example usage of the decorator with connection management
# @session_manager.connection(isolation_level="SERIALIZABLE", commit=True)
# async def example_method(*args, session: AsyncSession, **kwargs):
#     # Method logic
#     pass

# Example of dependency usage in a FastAPI route
# @router.post("/register/")
# async def register_user(user_data: SUserRegister, session: AsyncSession = TransactionSessionDep):
#     # Endpoint logic
#     pass
