CREATE TABLE IF NOT EXISTS tenants (
  tenant_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  phone TEXT NOT NULL,
  house_id UUID NOT NULL REFERENCES houses ON DELETE CASCADE,
  personal_id_type TEXT NOT NULL,
  personal_id TEXT NOT NULL,
  photo BYTEA ,
  active BOOL NOT NULL ,
  soc DATE NOT NULL,
  eos DATE NOT NULL

); 