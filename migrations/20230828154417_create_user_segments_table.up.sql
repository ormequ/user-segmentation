CREATE TABLE user_segments
(
    user_id    BIGINT NOT NULL,
    segment_id BIGINT NOT NULL
        CONSTRAINT user_segments_segment_id_key NOT NULL REFERENCES segments (id) ON DELETE CASCADE,
    CONSTRAINT user_segments_pkey PRIMARY KEY (user_id, segment_id)
);
