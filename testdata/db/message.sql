INSERT INTO channel(project_name, channel_name) VALUES ('scripbox', 'ops-incident');

INSERT INTO message(user, text, ts, channel_id, parent_id) VALUES ('mohit', 'hello world', '123', 1, -1);

INSERT INTO message(user, text, ts, channel_id, parent_id) VALUES ('bot', 'howdy', '124', 1, 1)
