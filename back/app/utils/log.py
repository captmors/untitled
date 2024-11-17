import logging
import sys
from loguru import logger

class InterceptHandler(logging.Handler):
    loglevel_mapping = {
        50: "CRITICAL",
        40: "ERROR",
        30: "WARNING",
        20: "INFO",
        10: "DEBUG",
        0: "NOTSET",
    }

    def emit(self, record):
        level = self.loglevel_mapping.get(record.levelno, record.levelno)
        logger.opt(depth=6, exception=record.exc_info).log(level, record.getMessage())

def setup_logging(log_file_path="logs/app.log"):
    logger.remove()
    
    fmt = "<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan> - <level>{message}</level>" 

    logger.add(
        log_file_path,
        rotation="10 MB",   
        retention="10 days",
        format=fmt,
        level="INFO",
        mode="w",
        backtrace=False,
        diagnose=False
    )
    
    logger.add(
        sys.stdout,
        format=fmt,
        level="INFO",
        backtrace=False,
        diagnose=False
    )

    logging.basicConfig(handlers=[InterceptHandler()], level=0)

    for log_name in ("uvicorn", "uvicorn.access", "uvicorn.error", "fastapi"):
        logging_logger = logging.getLogger(log_name)
        logging_logger.handlers = [InterceptHandler()]
        logging_logger.propagate = False

