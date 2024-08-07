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
  user_id UUID REFERENCES users (id),
  post_id UUID REFERENCES posts (id),
  UNIQUE (user_id, post_id)
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
  comment_id UUID REFERENCES comments (id) UNIQUE
);

CREATE TABLE comment_likes (
  user_id UUID REFERENCES users (id),
  comment_id UUID REFERENCES comments (id),
  UNIQUE (user_id, comment_id)
);

-- Data samples
INSERT INTO
  users (id, first_name, last_name, nick, passhash)
VALUES
  (
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'Ilya',
    'Blinkov',
    'drugtester',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  ),
  (
    '48683858-796c-45ad-a361-9e3d6d003354',
    'Marina',
    'Graf',
    'mmm',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  ),
  (
    'baa0e8bc-385f-4314-9580-29855aff2229',
    'Sasha',
    'Provodnikov',
    'yyy',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  ),
  (
    '37d28fdf-99bc-44b5-8df9-6a3b1a36f177',
    'Tolber',
    'Ovcharenko',
    'tolber01',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  ),
  (
    '2db640fd-7aa4-4bba-8ee6-3935b700297a',
    'Sanal',
    'Mandjiev',
    'sanal',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  );

INSERT INTO
  comments (id, author_id, message, created_at)
VALUES
  (
    '620c79b7-3927-48b7-a308-1ffd3db6036f',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'comment message',
    TIMESTAMP '2024-02-16 15:36:38'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661fa',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '2nd comment message',
    TIMESTAMP '2024-02-16 15:36:55'
  );

INSERT INTO
  posts (id, author_id, owner_id, message, created_at)
VALUES
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message',
    TIMESTAMP '2024-02-16 14:36:38'
  );

INSERT INTO
  post_comments (post_id, comment_id)
VALUES
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661fa'
  );
