ALTER TABLE users
ADD CONSTRAINT username_unique UNIQUE (username);

ALTER TABLE users
ADD CONSTRAINT phone_unique UNIQUE (phone);
