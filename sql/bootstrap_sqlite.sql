CREATE TABLE "users" (
    "id"            TEXT,
    "username"      TEXT NOT NULL,
    "access_token"  TEXT NOT NULL UNIQUE,
    "refresh_token" TEXT NOT NULL UNIQUE,
    "token_expiry"  INTEGER NOT NULL,

    CONSTRAINT "User_PK" PRIMARY KEY("id")
);

CREATE TABLE "anilist_users" (
    "id"            INTEGER,
    "access_token"  TEXT NOT NULL UNIQUE,
    "refresh_token" TEXT NOT NULL UNIQUE,
    "token_expiry"  INTEGER NOT NULL,
    "user_id"       TEXT NOT NULL,

    CONSTRAINT "AnilistUser_PK" PRIMARY KEY("id"),
    CONSTRAINT "User_FK"        FOREIGN KEY("user_id") REFERENCES "users"("id")
);
