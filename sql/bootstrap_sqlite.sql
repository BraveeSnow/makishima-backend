CREATE TABLE "users" (
	"id"	        TEXT,
	"username"	    TEXT NOT NULL,
	"access_token"	TEXT NOT NULL UNIQUE,
	"refresh_token"	TEXT NOT NULL UNIQUE,
	"token_expiry"	INTEGER NOT NULL,
	PRIMARY KEY("id")
)
