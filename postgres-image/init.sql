-- Definitions
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  first_name VARCHAR(64) NOT NULL,
  last_name VARCHAR(64) NOT NULL,
  nick VARCHAR(64) NOT NULL UNIQUE,
  passhash VARCHAR(64) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions (
  user_id UUID REFERENCES users (id) UNIQUE,
  session UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_activity_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE friendship_requests (
  requester_id UUID REFERENCES users (id),
  responder_id UUID REFERENCES users (id),
  message VARCHAR(256),
  accepted BOOLEAN NOT NULL DEFAULT FALSE,
  viewed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  author_id UUID REFERENCES users (id),
  owner_id UUID REFERENCES users (id),
  message VARCHAR(4096),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  removed_at TIMESTAMP
);

CREATE TABLE post_likes (
  author_id UUID REFERENCES users (id),
  post_id UUID REFERENCES posts (id)
);

CREATE TABLE comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  author_id UUID REFERENCES users (id),
  message VARCHAR(4096),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  removed_at TIMESTAMP
);

CREATE TABLE post_comments (
  post_id UUID REFERENCES posts (id),
  comment_id UUID REFERENCES comments (id)
);

CREATE TABLE comment_likes (
  author_id UUID REFERENCES users (id),
  comment_id UUID REFERENCES comments (id)
);

-- Data samples
INSERT INTO
  users (first_name, last_name, nick, passhash)
VALUES
  ('Ilya', 'Blinkov', 'zzz', 'password'),
  ('Marina', 'Graf', 'mmm', 'password'),
  ('Sasha', 'Provodnikov', 'yyy', 'password'),
  ('Tolber', 'Ovcharenko', 'tolber01', 'password'),
  ('Sanal', 'Mandjiev', 'sanal', 'password');
