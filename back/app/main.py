from fastapi import FastAPI
from app.utils.handlers import add_exception_handlers, add_cors_middleware
from app.utils.log import setup_logging
from loguru import logger
from app.auth.router import router as auth_router
from app.music.router import router as music_router

app = FastAPI()

app.include_router(auth_router)
app.include_router(music_router)

add_exception_handlers(app)
add_cors_middleware(app)

setup_logging()
