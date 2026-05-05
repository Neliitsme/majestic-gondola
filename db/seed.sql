-- Test seed data
-- Run with: psql $POSTGRES_URL -f db/seed.sql

BEGIN;

-- Artists
INSERT INTO "artists" ("artist_id", "name") OVERRIDING SYSTEM VALUE VALUES
  (1, 'The Beatles'),
  (2, 'Radiohead'),
  (3, 'Solo Artist');

-- Tracks
-- Track 2 and 4 already have processed reviews (score > 0, review_count > 0)
-- Track 1, 3, 5 start fresh (score=0, review_count=0) — unprocessed reviews target these
INSERT INTO "tracks" ("track_id", "name", "artist_id", "score", "review_count") OVERRIDING SYSTEM VALUE VALUES
  (1, 'Come Together',  1, 0,  0),
  (2, 'Hey Jude',       1, 60, 3),
  (3, 'Creep',          2, 0,  0),
  (4, 'Karma Police',   2, 70, 2),
  (5, 'Only Track',     3, 0,  0),
  (6, 'No Artist',   NULL, 0,  0);

-- Users
INSERT INTO "users" ("user_id", "name") OVERRIDING SYSTEM VALUE VALUES
  (1, 'Alice'),
  (2, 'Bob'),
  (3, 'Carol');

-- Unprocessed reviews
INSERT INTO "reviews" ("user_id", "track_id", "score", "is_processed") VALUES
  (1, 1, 80, false),  -- Alice reviews Come Together
  (2, 1, 60, false),  -- Bob   reviews Come Together
  (1, 3, 90, false),  -- Alice reviews Creep
  (3, 3, 70, false),  -- Carol reviews Creep
  (2, 5, 50, false),  -- Bob   reviews Only Track
  (3, 6, 75, false),  -- Carol reviews No Artist (track_id set, artist_id NULL on track)
  (1, NULL, 85, false); -- review with no track_id (should be skipped by processor)

COMMIT;

-- =============================================================================
-- Expected results after review processor runs
-- =============================================================================
--
-- Formula: newScore = (currentScore * currentCount + sumOfNewScores) / (currentCount + newCount)
--
-- Track 1 (Come Together): (0*0 + 80+60) / (0+2) = 140/2 = 70,  review_count = 2
-- Track 3 (Creep):         (0*0 + 90+70) / (0+2) = 160/2 = 80,  review_count = 2
-- Track 5 (Only Track):    (0*0 + 50)    / (0+1) = 50/1  = 50,  review_count = 1
-- Track 6 (No Artist):     (0*0 + 75)    / (0+1) = 75/1  = 75,  review_count = 1
--   (track 6 score updates even though its track has no artist)
--
-- All 6 targeted reviews marked is_processed = true.
-- Review with track_id NULL is skipped (stays unprocessed).
--
-- =============================================================================
-- Expected results after artist processor runs (uses updated track scores above)
-- =============================================================================
--
-- Artist 1 (The Beatles):
--   Track 1 score=70 rc=2, Track 2 score=60 rc=3 → avg = (70+60)/2 = 65
--
-- Artist 2 (Radiohead):
--   Track 3 score=80 rc=2, Track 4 score=70 rc=2 → avg = (80+70)/2 = 75
--
-- Artist 3 (Solo Artist):
--   Track 5 score=50 rc=1 → avg = 50
--
-- Track 6 has no artist_id so it doesn't affect any artist score.
