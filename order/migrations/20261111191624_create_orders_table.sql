-- +goose Up
create table if not exists orders (
    order_uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    part_uuids TEXT[] NOT NULL,
    total_price DOUBLE PRECISION NOT NULL,
    transaction_uuid UUID,
    payment_method TEXT CHECK (payment_method IS NULL OR payment_method IN ('UNKNOWN','CARD','SBP','CREDIT_CARD','INVESTOR_MONEY')),
    status TEXT NOT NULL CHECK (status in ('PENDING_PAYMENT','PAID','CANCELLED'))
);

-- +goose Down
drop table if exists orders;
