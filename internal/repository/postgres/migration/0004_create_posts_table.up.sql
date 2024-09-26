CREATE TABLE IF NOT EXISTS posts
(
    id        uuid PRIMARY KEY      DEFAULT uuid_generate_v4(),
    title     TEXT         NOT NULL,
    author    VARCHAR(255) NOT NULL,
    origin    VARCHAR(255) NOT NULL,
    synced_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata  JSONB        NOT NULL,

    UNIQUE (metadata)
);