CREATE TABLE IF NOT EXISTS "user" (
	"id" INTEGER NOT NULL,
	"username" VARCHAR NOT NULL,
	"password" VARCHAR NOT NULL,
	"last_login_ip" VARCHAR,
	"last_login_date" TIMESTAMP,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "news" (
	"id" INTEGER NOT NULL,
	"content" VARCHAR,
	"update_time" TIMESTAMP,
	PRIMARY KEY("id")
);
