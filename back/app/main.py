from fastapi import FastAPI
from app.utils.handlers import add_exception_handlers, add_cors_middleware
from app.utils.log import setup_logging
from loguru import logger

app = FastAPI()

add_exception_handlers(app)
add_cors_middleware(app)

setup_logging()

logger.info("hello loguru")
