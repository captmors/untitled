import os
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    BASE_DIR: str = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
    
    DB_URL: str = "postgresql+asyncpg://postgres:postgres@localhost:5432/postgres"
    DISABLE_AUTH: str = "true"

    model_config = SettingsConfigDict(env_file=f"{BASE_DIR}/.env")


settings = Settings()

DB_URL = settings.DB_URL
