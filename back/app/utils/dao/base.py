from typing import List, Any, TypeVar, Generic
from pydantic import BaseModel
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.future import select
from sqlalchemy import update as sqlalchemy_update, delete as sqlalchemy_delete, func
from loguru import logger
from sqlalchemy.ext.asyncio import AsyncSession
from .database import Base

# Define a type parameter T with a constraint that it inherits from Base
T = TypeVar("T", bound=Base)


class BaseDAO(Generic[T]):
    model: type[T]

    @classmethod
    async def find_one_or_none_by_id(cls, data_id: int, session: AsyncSession):
        # Find a record by its ID

        logger.info(f"Searching for {cls.model.__name__} with ID: {data_id}")
        try:
            query = select(cls.model).filter_by(id=data_id)
            result = await session.execute(query)
            record = result.scalar_one_or_none()
            if record:
                logger.info(f"Record with ID {data_id} found.")
            else:
                logger.info(f"Record with ID {data_id} not found.")
            return record
        except SQLAlchemyError as e:
            logger.error(f"Error finding record with ID {data_id}: {e}")
            raise

    @classmethod
    async def find_one_or_none(cls, session: AsyncSession, filters: BaseModel):
        # Find a single record by filters
        filter_dict = filters.model_dump(exclude_unset=True)
        logger.info(f"Searching for one {cls.model.__name__} record with filters: {filter_dict}")
        try:
            query = select(cls.model).filter_by(**filter_dict)
            result = await session.execute(query)
            record = result.scalar_one_or_none()
            if record:
                logger.info(f"Record found with filters: {filter_dict}")
            else:
                logger.info(f"No record found with filters: {filter_dict}")
            return record
        except SQLAlchemyError as e:
            logger.error(f"Error finding record with filters {filter_dict}: {e}")
            raise

    @classmethod
    async def find_all(cls, session: AsyncSession, filters: BaseModel | None):
        # Find all records matching filters
        filter_dict = filters.model_dump(exclude_unset=True) if filters else {}
        logger.info(f"Searching for all {cls.model.__name__} records with filters: {filter_dict}")
        try:
            query = select(cls.model).filter_by(**filter_dict)
            result = await session.execute(query)
            records = result.scalars().all()
            logger.info(f"Found {len(records)} records.")
            return records
        except SQLAlchemyError as e:
            logger.error(f"Error finding all records with filters {filter_dict}: {e}")
            raise

    @classmethod
    async def add(cls, session: AsyncSession, values: BaseModel):
        # Add a single record
        values_dict = values.model_dump(exclude_unset=True)
        logger.info(f"Adding {cls.model.__name__} record with parameters: {values_dict}")
        new_instance = cls.model(**values_dict)
        session.add(new_instance)
        try:
            await session.flush()
            logger.info(f"{cls.model.__name__} record added successfully.")
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error adding record: {e}")
            raise e
        return new_instance

    @classmethod
    async def add_many(cls, session: AsyncSession, instances: List[BaseModel]):
        # Add multiple records
        values_list = [item.model_dump(exclude_unset=True) for item in instances]
        logger.info(f"Adding multiple {cls.model.__name__} records. Count: {len(values_list)}")
        new_instances = [cls.model(**values) for values in values_list]
        session.add_all(new_instances)
        try:
            await session.flush()
            logger.info(f"Successfully added {len(new_instances)} records.")
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error adding multiple records: {e}")
            raise e
        return new_instances

    @classmethod
    async def update(cls, session: AsyncSession, filters: BaseModel, values: BaseModel):
        # Update records matching filters
        filter_dict = filters.model_dump(exclude_unset=True)
        values_dict = values.model_dump(exclude_unset=True)
        logger.info(f"Updating {cls.model.__name__} records with filter: {filter_dict} and values: {values_dict}")
        query = (
            sqlalchemy_update(cls.model)
            .where(*[getattr(cls.model, k) == v for k, v in filter_dict.items()])
            .values(**values_dict)
            .execution_options(synchronize_session="fetch")
        )
        try:
            result = await session.execute(query)
            await session.flush()
            logger.info(f"{result.rowcount} records updated.")
            return result.rowcount
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error updating records: {e}")
            raise e

    @classmethod
    async def delete(cls, session: AsyncSession, filters: BaseModel):
        # Delete records matching filters
        filter_dict = filters.model_dump(exclude_unset=True)
        logger.info(f"Deleting {cls.model.__name__} records with filter: {filter_dict}")
        if not filter_dict:
            logger.error("At least one filter is required for deletion.")
            raise ValueError("At least one filter is required for deletion.")

        query = sqlalchemy_delete(cls.model).filter_by(**filter_dict)
        try:
            result = await session.execute(query)
            await session.flush()
            logger.info(f"{result.rowcount} records deleted.")
            return result.rowcount
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error deleting records: {e}")
            raise e

    @classmethod
    async def count(cls, session: AsyncSession, filters: BaseModel):
        # Count the number of records matching filters
        filter_dict = filters.model_dump(exclude_unset=True)
        logger.info(f"Counting {cls.model.__name__} records with filter: {filter_dict}")
        try:
            query = select(func.count(cls.model.id)).filter_by(**filter_dict)
            result = await session.execute(query)
            count = result.scalar()
            logger.info(f"Found {count} records.")
            return count
        except SQLAlchemyError as e:
            logger.error(f"Error counting records: {e}")
            raise

    @classmethod
    async def paginate(cls, session: AsyncSession, page: int = 1, page_size: int = 10, filters: BaseModel = None):
        # Paginate records
        filter_dict = filters.model_dump(exclude_unset=True) if filters else {}
        logger.info(
            f"Paginating {cls.model.__name__} records with filter: {filter_dict}, page: {page}, page size: {page_size}")
        try:
            query = select(cls.model).filter_by(**filter_dict)
            result = await session.execute(query.offset((page - 1) * page_size).limit(page_size))
            records = result.scalars().all()
            logger.info(f"Found {len(records)} records on page {page}.")
            return records
        except SQLAlchemyError as e:
            logger.error(f"Error paginating records: {e}")
            raise

    @classmethod
    async def find_by_ids(cls, session: AsyncSession, ids: List[int]) -> List[Any]:
        """Find multiple records by a list of IDs"""
        logger.info(f"Finding {cls.model.__name__} records by ID list: {ids}")
        try:
            query = select(cls.model).filter(cls.model.id.in_(ids))
            result = await session.execute(query)
            records = result.scalars().all()
            logger.info(f"Found {len(records)} records by ID list.")
            return records
        except SQLAlchemyError as e:
            logger.error(f"Error finding records by ID list: {e}")
            raise

    @classmethod
    async def upsert(cls, session: AsyncSession, unique_fields: List[str], values: BaseModel):
        """Create or update a record"""
        values_dict = values.model_dump(exclude_unset=True)
        filter_dict = {field: values_dict[field] for field in unique_fields if field in values_dict}

        logger.info(f"Upsert for {cls.model.__name__}")
        try:
            existing = await cls.find_one_or_none(session, BaseModel.construct(**filter_dict))
            if existing:
                # Update the existing record
                for key, value in values_dict.items():
                    setattr(existing, key, value)
                await session.flush()
                logger.info(f"Updated existing {cls.model.__name__} record")
                return existing
            else:
                # Create a new record
                new_instance = cls.model(**values_dict)
                session.add(new_instance)
                await session.flush()
                logger.info(f"Created new {cls.model.__name__} record")
                return new_instance
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error in upsert: {e}")
            raise

    @classmethod
    async def bulk_update(cls, session: AsyncSession, records: List[BaseModel]) -> int:
        """Bulk update records"""
        logger.info(f"Bulk updating {cls.model.__name__} records")
        try:
            updated_count = 0
            for record in records:
                filter_dict = {"id": record.id}
                values_dict = record.model_dump(exclude_unset=True)
                query = (
                    sqlalchemy_update(cls.model)
                    .where(*[getattr(cls.model, k) == v for k, v in filter_dict.items()])
                    .values(**values_dict)
                    .execution_options(synchronize_session="fetch")
                )
                result = await session.execute(query)
                updated_count += result.rowcount
            await session.flush()
            logger.info(f"Bulk updated {updated_count} records.")
            return updated_count
        except SQLAlchemyError as e:
            await session.rollback()
            logger.error(f"Error in bulk update: {e}")
            raise
