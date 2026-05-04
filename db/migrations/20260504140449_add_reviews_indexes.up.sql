CREATE INDEX IF NOT EXISTS "ix_reviews_unprocessed" ON "reviews" ("review_id")
	INCLUDE ("track_id", "score")
	WHERE is_processed = false;

CREATE INDEX IF NOT EXISTS "ix_reviews_deleted" ON "reviews" ("review_id")
	INCLUDE ("track_id", "score")
	WHERE is_deleted = true;
