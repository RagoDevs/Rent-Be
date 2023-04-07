CREATE TABLE IF NOT EXISTS houses (
  house_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  location citext NOT NULL,
  block citext NOT NULL,
  partition SMALLINT NOT NULL,
  occupied BOOL NOT NULL
); 