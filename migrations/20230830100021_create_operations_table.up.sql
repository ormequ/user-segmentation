CREATE TABLE operations (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    segment_id BIGINT CONSTRAINT operations_segment_id_key NOT NULL REFERENCES segments (id) ON DELETE CASCADE,
    type SMALLINT NOT NULL,
    time TIMESTAMP NOT NULL
);