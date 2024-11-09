CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "code" varchar NOT NULL,
  "username" varchar NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" VARCHAR(20) NOT NULL,
  "description" varchar  NULL,
  "hashed_password" varchar NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz,
  "deleted" bool NOT NULL DEFAULT false
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" bigserial NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE TABLE role (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NULL,
  "deleted" bool NOT NULL DEFAULT false
);

CREATE TABLE permission_group (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL
);

CREATE TABLE permission (
    "id" bigserial PRIMARY KEY,
    "code" varchar NOT NULL,
    "name" varchar NOT NULL,
    "permission_group_id" bigint NOT NULL
);

ALTER TABLE "permission"
ADD CONSTRAINT fk_permission_permission_group FOREIGN KEY (permission_group_id) REFERENCES "permission_group" (id);


CREATE TABLE role_permission (
    "id" bigserial PRIMARY KEY,
    "role_id" bigint NOT NULL,
    "permission_id" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NULL,
    "deleted" bool NOT NULL DEFAULT false
);

ALTER TABLE "role_permission"
ADD CONSTRAINT fk_role_permission_role FOREIGN KEY (role_id) REFERENCES "role" (id);

ALTER TABLE "role_permission"
ADD CONSTRAINT fk_role_permission_permission FOREIGN KEY (permission_id) REFERENCES "permission" (id);

CREATE TABLE user_role (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "role_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NULL,
  "deleted" bool NOT NULL DEFAULT false
);

ALTER TABLE "user_role"
ADD CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES "users" (id);

ALTER TABLE "user_role"
ADD CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES "role" (id);
