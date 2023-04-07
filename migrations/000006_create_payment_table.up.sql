CREATE TABLE IF NOT EXISTS payments (
    payment_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    tenant_id  UUID NOT NULL REFERENCES tenants ON DELETE CASCADE,
    period int NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);
