CREATE TYPE "task_status" AS ENUM (
  'pending',
  'done',
  'deleted',
  'suspended',
  'recurring'
);

CREATE TYPE "task_prior" AS ENUM (
  'A',
  'B',
  'C'
);

CREATE TABLE "role" (
  "id" varchar(64),
  "name" varchar(64)
);

CREATE TABLE "user" (
  "uuid" uuid PRIMARY KEY,
  "first_name" varchar(64),
  "last_name" varchar(64),
  "email" varchar(64),
  "password" varchar(256),
  "signup_at" timestampz,
  "last_login" timestampz
);

CREATE TABLE "project" (
  "uuid" uuid PRIMARY KEY,
  "owner_uuid" uuid,
  "title" varchar(64),
  "password" varchar(256)
);

CREATE TABLE "event" (
  "uuid" uuid PRIMARY KEY,
  "title" varchar(64),
  "created_at" timestampz,
  "updated_at" timestampz,
  "start_at" timestampz,
  "end_at" timestampz,
  "actual_end" timestampz,
  "description" text
);

CREATE TABLE "project_member" (
  "project_uuid" uuid,
  "member_uuid" uuid,
  "role" varchar(64),
  PRIMARY KEY ("project_uuid", "member_uuid")
);

CREATE TABLE "project_event" (
  "project_uuid" uuid,
  "event_uuid" uuid,
  PRIMARY KEY ("project_uuid", "event_uuid")
);

CREATE TABLE "timesheet_report" (
  "uuid" uuid PRIMARY KEY,
  "project_uuid" uuid,
  "member_uuid" uuid,
  "timesheet_uuid" uuid,
  "feedback_by" uuid,
  "title" varchar(64),
  "description" text,
  "feedback" text
);

CREATE TABLE "task_report" (
  "uuid" uuid PRIMARY KEY,
  "project_uuid" uuid,
  "member_uuid" uuid,
  "task_uuid" uuid,
  "feedback_by" uuid,
  "title" varchar(64),
  "description" text,
  "feedback" text
);

CREATE TABLE "task" (
  "uuid" uuid PRIMARY KEY,
  "owner_producer_uuid" uuid,
  "project_uuid" uuid,
  "member_uuid" uuid,
  "parent_uuid" uuid,
  "title" varchar(64),
  "description" text,
  "created_at" timestampz,
  "updated_at" timestampz
);

CREATE TABLE "task_detail" (
  "task_uuid" uuid,
  "id" integer,
  "end_at" timestampz,
  "deadline" timestampz,
  "schedule" timestampz,
  "tags" varchar[],
  "status" task_status,
  "priority" task_prior,
  PRIMARY KEY ("task_uuid")
);

CREATE TABLE "timesheet" (
  "uuid" uuid PRIMARY KEY,
  "user_uuid" uuid,
  "created_at" timestampz,
  "updated_at" timestamp
);

CREATE TABLE "task_time" (
  "uuid" uuid PRIMARY KEY,
  "task_uuid" uuid,
  "timesheet_uuid" uuid,
  "created_at" timestampz,
  "tracked_time" interval,
  "billable" float
);

CREATE TABLE "calendar" (
  "uuid" uuid PRIMARY KEY,
  "user_uuid" uuid,
  "task_uuid" uuid,
  "created_at" timestampz,
  "updated_at" timestampz,
  "started_at" timestampz,
  "duration" interval
);

ALTER TABLE "project" ADD FOREIGN KEY ("owner_uuid") REFERENCES "user" ("uuid");

ALTER TABLE "project_member" ADD FOREIGN KEY ("project_uuid") REFERENCES "project" ("uuid");

ALTER TABLE "project_member" ADD FOREIGN KEY ("member_uuid") REFERENCES "user" ("uuid");

ALTER TABLE "project_member" ADD FOREIGN KEY ("role") REFERENCES "role" ("id");

ALTER TABLE "project_event" ADD FOREIGN KEY ("project_uuid") REFERENCES "project" ("uuid");

ALTER TABLE "project_event" ADD FOREIGN KEY ("event_uuid") REFERENCES "event" ("uuid");

ALTER TABLE "timesheet_report" ADD FOREIGN KEY ("timesheet_uuid") REFERENCES "timesheet" ("uuid");

ALTER TABLE "timesheet_report" ADD FOREIGN KEY ("feedback_by") REFERENCES "user" ("uuid");

ALTER TABLE "timesheet_report" ADD FOREIGN KEY ("project_uuid", "member_uuid") REFERENCES "project_member" ("project_uuid", "member_uuid");

ALTER TABLE "task_report" ADD FOREIGN KEY ("task_uuid") REFERENCES "task" ("uuid");

ALTER TABLE "task_report" ADD FOREIGN KEY ("feedback_by") REFERENCES "user" ("uuid");

ALTER TABLE "task_report" ADD FOREIGN KEY ("project_uuid", "member_uuid") REFERENCES "project_member" ("project_uuid", "member_uuid");

ALTER TABLE "task" ADD FOREIGN KEY ("owner_producer_uuid") REFERENCES "user" ("uuid");

ALTER TABLE "task" ADD FOREIGN KEY ("parent_uuid") REFERENCES "task" ("uuid");

ALTER TABLE "task" ADD FOREIGN KEY ("project_uuid", "member_uuid") REFERENCES "project_member" ("project_uuid", "member_uuid");

ALTER TABLE "task_detail" ADD FOREIGN KEY ("task_uuid") REFERENCES "task" ("uuid");

ALTER TABLE "timesheet" ADD FOREIGN KEY ("user_uuid") REFERENCES "user" ("uuid");

ALTER TABLE "task_time" ADD FOREIGN KEY ("task_uuid") REFERENCES "task" ("uuid");

ALTER TABLE "task_time" ADD FOREIGN KEY ("timesheet_uuid") REFERENCES "timesheet" ("uuid");

ALTER TABLE "calendar" ADD FOREIGN KEY ("user_uuid") REFERENCES "user" ("uuid");

ALTER TABLE "calendar" ADD FOREIGN KEY ("task_uuid") REFERENCES "task" ("uuid");

ALTER TABLE "timesheet" ADD FOREIGN KEY ("created_at") REFERENCES "timesheet" ("user_uuid");
