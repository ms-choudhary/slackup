CREATE TABLE channel (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  channel_name TEXT,
  project_name TEXT
);

CREATE TABLE message (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user TEXT,
  text TEXT,
  ts TEXT,
  channel_id INTEGER,
  parent_id INTEGER,
  FOREIGN KEY(channel_id) REFERENCES channel(id)
);
