CREATE TABLE IF NOT EXISTS tokens (
hash bytea PRIMARY KEY,
admin_id UUID NOT NULL REFERENCES admins ON DELETE CASCADE,
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);