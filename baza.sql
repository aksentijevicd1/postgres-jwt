CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "firstname" varchar NOT NULL,
    "lastname" varchar NOT NULL,
    "password" varchar NOT NULL,
    "email" varchar NOT NULL,
    "phone" varchar NOT NULL,
    "token" varchar NOT NULL,
    "user_type_id" bigint NOT NULL,
    "refresh_token" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT 'now()',
    "updated_at" timestamptz NOT NULL DEFAULT 'now()',
);

CREATE TABLE "usertypes" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL
);

CREATE TABLE "tasks" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "illustration" varchar NOT NULL,
  "category_id" bigint NOT NULL,
  "status_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "statuses" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("phone");

CREATE INDEX ON "tasks" ("title");

CREATE INDEX ON "tasks" ("category_id");

CREATE INDEX ON "tasks" ("status_id");

CREATE INDEX ON "categories" ("name");

CREATE INDEX ON "statuses" ("name");

ALTER TABLE "tasks" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("status_id") REFERENCES "statuses" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("user_type_id") REFERENCES "usertypes" ("id");
