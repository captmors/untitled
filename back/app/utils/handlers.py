from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from app.auth.utils import create_admin_role_and_user
from app.utils.dao.session_maker import session_manager
from app.utils.s3.routine import minio_routine
from loguru import logger 

async def lifespan(app: FastAPI):
    logger.info("Database routines...")
    async with session_manager.create_session() as session:
        await create_admin_role_and_user(session)
    
    logger.info("Minio routine...")        
    minio_routine()
    
    yield

def add_exception_handlers(app: FastAPI):
    @app.exception_handler(HTTPException)
    async def http_exception_handler(request: Request, exc: HTTPException):
        return JSONResponse(
            status_code=exc.status_code,
            content={"detail": exc.detail},
        )

    @app.exception_handler(Exception)
    async def general_exception_handler(request: Request, exc: Exception):
        return JSONResponse(
            status_code=500,
            content={"message": "An error occurred", "details": str(exc)},
        )

def add_cors_middleware(app: FastAPI):
    origins = [
        "http://localhost:3000",
        "http://localhost:5173",
        ]
    
    app.add_middleware(
        CORSMiddleware,
        allow_origins=origins,  
        allow_credentials=True,
        allow_methods=["*"],  
        allow_headers=["*"], 
    )
    