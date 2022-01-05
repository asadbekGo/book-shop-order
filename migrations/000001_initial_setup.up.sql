CREATE TABLE IF NOT EXISTS orders(
    order_id uuid not null primary key,
    book_id uuid not null,
    quantity int not null,
    description text,
    created_at timestamp default current_timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
