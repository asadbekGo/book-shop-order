ALTER TABLE order ADD COLUMN created_at timestamp default current_timestamp;
ALTER TABLE order ADD COLUMN updated_at timestamp;
ALTER TABLE order ADD COLUMN deleted_at timestamp;
