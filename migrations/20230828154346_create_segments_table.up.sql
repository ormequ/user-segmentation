CREATE TABLE segments
(
    id   SERIAL PRIMARY KEY,
    slug VARCHAR(255) CONSTRAINT segments_slug_key UNIQUE NOT NULL
);
CREATE INDEX ON segments(slug);
