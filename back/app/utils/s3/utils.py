import uuid
from minio import Minio

from .routine import minio_client, USER_AVATARS_BUCKET

def upload_file_to_minio(bucket_name: str, file_data: bytes, file_name: str) -> str:
    """
    Загружает файл в указанный бакет.
    Возвращает URL файла.
    """
    unique_name = f"{uuid.uuid4()}_{file_name}"
    minio_client.put_object(bucket_name, unique_name, file_data, len(file_data))
    return f"http://{minio_client._endpoint}/{bucket_name}/{unique_name}"

def get_file_from_minio(bucket_name: str, file_name: str) -> bytes:
    """
    Загружает файл из MinIO.
    """
    response = minio_client.get_object(bucket_name, file_name)
    return response.read()
