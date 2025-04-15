CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION citext;

CREATE TABLE IF NOT EXISTS admin (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    email citext UNIQUE NOT NULL,
    password_hash BYTEA NOT NULL,
    activated BOOL NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
);

CREATE TABLE IF NOT EXISTS token (
    hash BYTEA PRIMARY KEY,
    id UUID NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
    
);

CREATE TABLE IF NOT EXISTS house (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    location CITEXT NOT NULL,
    block CITEXT NOT NULL,
    partition SMALLINT NOT NULL,
    occupied BOOL NOT NULL,
    price INT NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()
); 


CREATE TABLE IF NOT EXISTS tenant(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    house_id UUID NOT NULL REFERENCES house(id) ON DELETE CASCADE,
    personal_id_type TEXT NOT NULL DEFAULT '',
    personal_id TEXT NOT NULL DEFAULT '',
    photo TEXT NOT NULL DEFAULT '',
    active BOOL NOT NULL ,
    sos DATE NOT NULL,
    eos DATE NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4()

); 

CREATE TABLE IF NOT EXISTS payment (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    tenant_id  UUID NOT NULL REFERENCES tenant(id) ON DELETE CASCADE,
    amount INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    version UUID NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
    
);


