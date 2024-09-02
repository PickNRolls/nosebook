-- Definitions
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  first_name VARCHAR(64) NOT NULL,
  last_name VARCHAR(64) NOT NULL,
  nick VARCHAR(64) NOT NULL UNIQUE,
  passhash VARCHAR(64) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_activity_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions (
  user_id UUID REFERENCES users (id),
  session_id UUID NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL DEFAULT NOW() + '1day'::interval
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
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'Test',
    'Tester',
    'test_tester',
    '$2a$04$PFIkrnjZ62TLHhcU3a6Breh1sLUVMXzwlrLNo2dqQSTM9d02py.oa' -- 123123123 unhashed
  ),
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
  ),
  (
    '12d28fdf-99bc-44b5-8df9-6a3b1a36f177',
    'Somebody',
    'Some',
    'somebody',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  ),
  (
    '78b640fd-7aa4-4bba-8ee6-3935b700297a',
    'john',
    'wick',
    'killer',
    '$2a$04$A08tmv8hEQkc75GbpRlpMO6ClwAwEfEO0I1YG2qB56o/jsOdtn3hS'
  );

INSERT INTO
  user_sessions (user_id, session_id, expires_at)
VALUES
  (
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'bb23af03-be50-4bce-b729-b259b2e02e54',
    TIMESTAMP '2050-02-16 15:36:55'
  ),
  (
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'bb23af03-be50-4bce-b729-b259b2e02e55',
    TIMESTAMP '2050-02-16 15:36:55'
  );

INSERT INTO
  friendship_requests (requester_id, responder_id, message, created_at)
VALUES
  (
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'test request',
    TIMESTAMP '2024-02-16 15:36:38'
  ),
  (
    'baa0e8bc-385f-4314-9580-29855aff2229',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'test request',
    TIMESTAMP '2024-02-16 15:36:38'
  );

INSERT INTO
  friendship_requests (requester_id, responder_id, message, accepted, viewed, created_at)
VALUES
  (
    '37d28fdf-99bc-44b5-8df9-6a3b1a36f177',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'test',
    true,
    true,
    TIMESTAMP '2024-02-16 15:36:38'
  ),
  (
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    '2db640fd-7aa4-4bba-8ee6-3935b700297a',
    'test',
    true,
    true,
    TIMESTAMP '2024-02-16 15:36:38'
  );

INSERT INTO
  comments (id, author_id, message, created_at)
VALUES
  (
    '120c79b8-3927-48b7-a308-1ffd3db6036f',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'comment for permissions test',
    TIMESTAMP '2024-02-16 15:36:38'
  ),
  (
    '620c79b7-3927-48b7-a308-1ffd3db6036f',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'comment message',
    TIMESTAMP '2024-02-16 15:36:38'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f0',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '3nd comment message',
    TIMESTAMP '2024-02-16 15:36:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f1',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '4nd comment message',
    TIMESTAMP '2024-02-16 15:37:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f2',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '5nd comment message',
    TIMESTAMP '2024-02-16 15:38:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f3',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '6nd comment message',
    TIMESTAMP '2024-02-16 15:39:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '7nd comment message',
    TIMESTAMP '2024-02-16 15:40:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f5',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '8nd comment message',
    TIMESTAMP '2024-02-16 15:41:55'
  ),
  (
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f6',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '9nd comment message',
    TIMESTAMP '2024-02-16 15:42:55'
  );


INSERT INTO
  posts (id, author_id, owner_id, message, created_at)
VALUES
  (
    '27b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'test permissions post with only author id',
    TIMESTAMP '2024-02-16 14:36:38'
  ),
  (
    '27b7bf27-38f9-4ed5-b0a8-501a90f7c8e7',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'ed1a3fd0-4d0b-4961-b4cd-cf212357740d',
    'test permissions post with only owner id',
    TIMESTAMP '2024-02-16 14:36:38'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message',
    TIMESTAMP '2024-02-16 14:36:38'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e8',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 2',
    TIMESTAMP '2024-02-16 14:36:48'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e9',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 3',
    TIMESTAMP '2024-02-16 14:36:53'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e0',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 4',
    TIMESTAMP '2024-02-16 14:36:57'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e2',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 5',
    TIMESTAMP '2024-02-16 14:36:59'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-601a90f7c8e7',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 6',
    TIMESTAMP '2024-02-16 14:36:38'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-701a90f7c8e8',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 7',
    TIMESTAMP '2024-02-16 14:36:48'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-801a90f7c8e9',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 8',
    TIMESTAMP '2024-02-16 14:36:53'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-901a90f7c8e0',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 9',
    TIMESTAMP '2024-02-16 14:36:57'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-001a90f7c8e2',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 10',
    TIMESTAMP '2024-02-16 14:36:59'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-011a90f7c8e2',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 11',
    TIMESTAMP '2024-02-16 14:36:59'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c807',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    '48683858-796c-45ad-a361-9e3d6d003354',
    'post message',
    TIMESTAMP '2024-02-16 14:36:38'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c818',
    '48683858-796c-45ad-a361-9e3d6d003354',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 2',
    TIMESTAMP '2024-02-16 14:36:48'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c829',
    '48683858-796c-45ad-a361-9e3d6d003354',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 3',
    TIMESTAMP '2024-02-16 14:36:53'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c830',
    '48683858-796c-45ad-a361-9e3d6d003354',
    '48683858-796c-45ad-a361-9e3d6d003354',
    'post message 4',
    TIMESTAMP '2024-02-16 14:36:57'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c842',
    '48683858-796c-45ad-a361-9e3d6d003354',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4',
    'post message 5',
    TIMESTAMP '2024-02-16 14:36:59'
  );

INSERT INTO
  post_likes (post_id, user_id)
VALUES
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e8',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e8',
    '48683858-796c-45ad-a361-9e3d6d003354'
  ),
  -- post 1

  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '1ae02f69-ea1a-4308-b825-0e5896e652e4'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '48683858-796c-45ad-a361-9e3d6d003354'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'baa0e8bc-385f-4314-9580-29855aff2229'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '37d28fdf-99bc-44b5-8df9-6a3b1a36f177'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '2db640fd-7aa4-4bba-8ee6-3935b700297a'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '12d28fdf-99bc-44b5-8df9-6a3b1a36f177'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '78b640fd-7aa4-4bba-8ee6-3935b700297a'
  );

INSERT INTO
  post_comments (post_id, comment_id)
VALUES
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    '120c79b8-3927-48b7-a308-1ffd3db6036f'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f0'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f1'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f2'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f3'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f4'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f5'
  ),
  (
    'c7b7bf17-38f9-4ed5-b0a8-501a90f7c8e7',
    'd0023f4d-8d7f-4907-9438-d2ed2a9661f6'
  );

