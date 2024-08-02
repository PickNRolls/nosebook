-- Definitions
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(64),
  last_name VARCHAR(64)
);

CREATE TABLE friendship_requests (
  requester_id INT REFERENCES users (id),
  responder_id INT REFERENCES users (id),
  accepted BOOLEAN,
  viewed BOOLEAN
);

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  author_id INTEGER REFERENCES users (id),
  message VARCHAR(4096)
);

CREATE TABLE post_likes (
  author_id INTEGER REFERENCES users (id),
  post_id INTEGER REFERENCES posts (id)
);

CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  author_id INTEGER REFERENCES users (id)
);

CREATE TABLE comment_likes (
  author_id INTEGER REFERENCES users (id),
  comment_id INTEGER REFERENCES comments (id)
);

-- Data samples
INSERT INTO
  users (first_name, last_name)
VALUES
  ('Ilya', 'Blinkov'),
  ('Marina', 'Graf'),
  ('Sasha', 'Provodnikov'),
  ('Tolber', 'Ovcharenko'),
  ('Sanal', 'Mandjiev');
