CREATE TABLE IF NOT EXISTS payments (
    payment_id UUID PRIMARY KEY NOT NULL,
    tenant_id  UUID NOT NULL REFERENCES tenants ON DELETE CASCADE,
    period int NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
