CREATE TABLE IF NOT EXISTS tokens (
hash bytea PRIMARY KEY,
token_id UUID NOT NULL REFERENCES admins ON DELETE CASCADE,
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);