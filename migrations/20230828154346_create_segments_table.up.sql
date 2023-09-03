CREATE TABLE segments
(
    id   BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    slug VARCHAR(255) CONSTRAINT segments_slug_key UNIQUE NOT NULL
);
CREATE INDEX ON segments(slug);
