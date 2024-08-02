CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  firstName VARCHAR(64),
  lastName VARCHAR(64)
);

CREATE TABLE friendship_requests (
  requesterId INT REFERENCES users (id),
  responderId INT REFERENCES users (id),
  accepted BOOLEAN,
  viewed BOOLEAN
);

CREATE TABLE posts (
  id SERIAL PRIMARY KEY,
  authorId INTEGER REFERENCES users (id),
  message VARCHAR(4096)
);

CREATE TABLE post_likes (
  authorId INTEGER REFERENCES users (id),
  postId INTEGER REFERENCES posts (id)
);

CREATE TABLE comments (
  id SERIAL PRIMARY KEY,
  authorId INTEGER REFERENCES users (id)
);

CREATE TABLE comment_likes (
  authorId INTEGER REFERENCES users (id),
  commentId INTEGER REFERENCES comments (id)
);
