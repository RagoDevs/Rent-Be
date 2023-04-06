CREATE TABLE IF NOT EXISTS houses (
  house_id UUID PRIMARY KEY NOT NULL,
  location TEXT NOT NULL,
  block TEXT NOT NULL,
  partition TEXT NOT NULL,
  occupied BOOL NOT NULL
); 