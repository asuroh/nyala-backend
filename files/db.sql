
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS "public"."roles";
CREATE TABLE "public"."roles" (
  "id" char(36) DEFAULT uuid_generate_v4 () NOT NULL,
  "data" jsonb NOT NULL,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6)
);
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" char(36) DEFAULT uuid_generate_v4 () NOT NULL,
  "data" jsonb NOT NULL,
  "role_id" char(36) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6)
);

ALTER TABLE "public"."roles" ADD CONSTRAINT "roles_pkey" PRIMARY KEY ("id");
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");
ALTER TABLE "public"."users" ADD CONSTRAINT "users_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON DELETE CASCADE ON UPDATE CASCADE;


BEGIN;
INSERT INTO "public"."roles" VALUES ('d57bfbfe-4979-4809-a151-f6cd30de657b', '{"role_name": "Member", "description": "Default role for register user"}', '2020-02-17 14:41:11.322647', '2020-02-17 14:41:11.322647', NULL);
INSERT INTO "public"."roles" VALUES ('381b7700-fd23-44b7-9d1f-befba9fa7d6a', '{"role_name": "Admin", "description": "Administrator"}', '2020-02-17 14:41:27.17323', '2020-02-17 14:41:27.17323', NULL);
COMMIT;

BEGIN;
INSERT INTO "public"."users" VALUES ('81b6e0e4-8be0-4656-aecf-e18a98c3a0a7', '{"email": "superadmin2@init.com", "status": {"is_active": true}, "password": "$2a$14$54ISc3vDoIG8bXr.S1TTxObKBiZ5dl9XDug2Grw1hRBTA2e5IEV4G", "username": "Superadmin 2"}', '381b7700-fd23-44b7-9d1f-befba9fa7d6a', '2020-11-23 02:48:55.863911', '2020-11-23 03:25:12.429625', NULL);
INSERT INTO "public"."users" VALUES ('fff76956-5cf6-4b2a-b571-9e078fa31fbc', '{"email": "admin@test.com", "status": {"is_active": true}, "password": "$2a$14$MtGxJqQsXyGjggX8Q2hDpOZn85wI3FWCiw.R0mNcxe20Kz9phCbW2", "username": "admin"}', 'd57bfbfe-4979-4809-a151-f6cd30de657b', '2020-11-23 03:26:15.229383', '2020-11-23 03:34:20.921368', '2020-11-23 03:34:20.921368');
COMMIT;
