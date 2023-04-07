CREATE TABLE IF NOT EXISTS houses (
  house_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  location TEXT NOT NULL,
  block TEXT NOT NULL,
  partition TEXT NOT NULL,
  occupied BOOL NOT NULL
); 