from minio import Minio

from app.utils.s3.config import MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_URL, USER_AVATARS_BUCKET

# Инициализация клиента MinIO
minio_client = Minio(
    MINIO_URL,
    access_key=MINIO_ACCESS_KEY,
    secret_key=MINIO_SECRET_KEY,
    secure=False
)

def initialize_buckets():
    """
    Создает бакеты, если их нет.
    """
    try:
        if not minio_client.bucket_exists(USER_AVATARS_BUCKET):
            minio_client.make_bucket(USER_AVATARS_BUCKET)
    except Exception as err:
        print(f"Ошибка создания бакета: {err}")

def minio_routine():
    initialize_buckets()
    