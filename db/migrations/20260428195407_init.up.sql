CREATE TABLE IF NOT EXISTS "tracks" (
	"track_id" integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" text NOT NULL,
	"artist_id" integer,
	"genres" text[],
	"release_date" timestamptz NOT NULL DEFAULT now(),
	"created_at" timestamptz NOT NULL DEFAULT now(),
	"score" integer NOT NULL DEFAULT 0,
	CONSTRAINT "tracks_pkey" PRIMARY KEY("track_id"),
	CONSTRAINT "chk_tracks_score" CHECK(score between 0 and 100)
);

CREATE TABLE IF NOT EXISTS "artists" (
	"artist_id" integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" text NOT NULL,
	"created_at" timestamptz NOT NULL DEFAULT now(),
	"score" integer NOT NULL DEFAULT 0,
	CONSTRAINT "pk_artists" PRIMARY KEY("artist_id"),
	CONSTRAINT "chk_artists_score" CHECK(score between 0 and 100)
);

CREATE TABLE IF NOT EXISTS "users" (
	"user_id" integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" text NOT NULL,
	CONSTRAINT "pk_users" PRIMARY KEY("user_id")
);

CREATE TABLE IF NOT EXISTS "reviews" (
	"review_id" integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	"user_id" integer,
	"is_deleted" boolean NOT NULL DEFAULT false,
	"created_at" timestamptz NOT NULL DEFAULT now(),
	"is_processed" boolean NOT NULL DEFAULT false,
	"score" integer NOT NULL DEFAULT 0,
	"track_id" integer,
	CONSTRAINT "pk_reviews" PRIMARY KEY("review_id"),
	CONSTRAINT "chk_reviews_score" CHECK(score between 0 and 100)
);

CREATE INDEX IF NOT EXISTS "ix_artists_name" ON "artists" (
	"name"
);

CREATE INDEX IF NOT EXISTS "ix_reviews_user_id" ON "reviews" (
	"user_id"
);

CREATE INDEX IF NOT EXISTS "ix_reviews_track_id" ON "reviews" (
	"track_id"
);

CREATE INDEX IF NOT EXISTS "ix_tracks_artist_id" ON "tracks" (
	"artist_id"
);

ALTER TABLE "tracks" ADD CONSTRAINT "fk_tracks_artists" FOREIGN KEY ("artist_id")
	REFERENCES "artists"("artist_id")
	ON DELETE RESTRICT
	ON UPDATE CASCADE
	NOT DEFERRABLE;

ALTER TABLE "reviews" ADD CONSTRAINT "fk_reviews_user_id" FOREIGN KEY ("user_id")
	REFERENCES "users"("user_id")
	ON DELETE SET NULL
	ON UPDATE CASCADE
	NOT DEFERRABLE;

ALTER TABLE "reviews" ADD CONSTRAINT "fk_reviews_track_id" FOREIGN KEY ("track_id")
	REFERENCES "tracks"("track_id")
	ON DELETE CASCADE
	ON UPDATE CASCADE
	NOT DEFERRABLE;

