CREATE TABLE IF NOT EXISTS messages
(
    uuid       UUID PRIMARY KEY,
    room       TEXT      NOT NULL,
    username   TEXT      NOT NULL,
    text       TEXT      NOT NULL,
    type       TEXT      NOT NULL DEFAULT 'message',
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_room_created
    ON messages (room, created_at DESC);
