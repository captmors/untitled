from fastapi import FastAPI
from app.utils.handlers import add_exception_handlers, add_cors_middleware, lifespan
from app.utils.log import setup_logging
from app.auth.api.auth import router as auth_router
from app.music.router import router as music_router

setup_logging()

app = FastAPI(lifespan=lifespan)

add_exception_handlers(app)
add_cors_middleware(app)

app.include_router(auth_router)
app.include_router(music_router)
